package main

//go:generate go run tools/genroutes.go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"

	"labdb.org/labdb/auth"
	"labdb.org/labdb/env"
	"labdb.org/labdb/models"
	"labdb.org/labdb/routes"
	"labdb.org/labdb/search"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var devProxyTarget = "http://localhost:3001"
var proxySuffix = "-backend.labdb.io"

func backendHost(c *gin.Context) string {
	target := devProxyTarget
	if env.Dev && os.Getenv("PROXY_TARGET") != "" {
		target = os.Getenv("PROXY_TARGET")
	}
	if env.Prod {
		target = "https://" + strings.Replace(c.Request.Host, ".labdb.io", proxySuffix, -1)
	}
	return target
}

func setupBackendRequest(req *http.Request, maybeCurrentUser string) {
	for k := range req.Header {
		if strings.HasPrefix(strings.ToLower(k), "cf-") {
			req.Header.Del(k)
		}
	}
	req.Header.Del("X-Forwarded-For")
	req.Header.Add("X-Labdb-Forwarded", "true")
	if maybeCurrentUser != "" {
		auth.AddAuthHeaders(maybeCurrentUser, req.Header)
	}
	// Unset the host, as the proxy sets the host using URL.Host, but the value on
	// the request itself will override it if present.
	req.Host = ""
}

func proxy(c *gin.Context) {
	url, err := url.Parse(backendHost(c))
	if err != nil {
		panic(err)
	}

	if c.Request.Header.Get("X-Labdb-Forwarded") == "true" {
		c.String(400, "Stuck in a recursive proxy loop.")
		return
	}
	setupBackendRequest(c.Request, auth.CurrentUserID(c))
	p := httputil.NewSingleHostReverseProxy(url)
	p.ServeHTTP(c.Writer, c.Request)
}

func redirectHTTPS(c *gin.Context) {
	usesTLSOnHeroku := c.Request.Header.Get("X-Forwarded-Proto") == "https"
	if env.Prod && !usesTLSOnHeroku {
		newTarget := fmt.Sprintf("https://%s%s", c.Request.Host, c.Request.RequestURI)
		c.Redirect(302, newTarget)
		c.Abort()
		return
	}
	c.Next()
}

func requireAuthorization(c *gin.Context) {
	u := auth.CurrentUser(c)
	if u.AuthWrite || (u.AuthRead && c.Request.Method == "GET") {
		c.Next()
		return
	}

	fmt.Printf("Access denied to %+v.\n", u)
	c.String(403, "Forbidden")
	c.Abort()
}

func startup() {
	if env.Prod {
		gin.SetMode(gin.ReleaseMode)
	}
}

func shutdown() {
	models.Shutdown()
}

func modelAPI(r *gin.Engine) {
	apiM := r.Group("/api/v1/m")

	apiM.GET("/:model/:id", func(c *gin.Context) {
		modelType := c.Param("model")
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(400, "Bad ID")
			c.Abort()
			return
		}
		if models.IsImplemented(modelType) {
			m := models.Empty(modelType)
			models.GetByID(m, id)
			if m.GetID() == 0 {
				c.String(404, "Not found.")
				c.Abort()
				return
			}
			c.JSON(200, m.AsResourceDef())
		} else {
			proxy(c)
		}
	})

	apiM.POST("/:model/new", func(c *gin.Context) {
		modelType := c.Param("model")
		if models.IsImplemented(modelType) {
			m := models.Empty(modelType)
			m.AutoFill(auth.CurrentUser(c).Name)
			models.Create(m)
			c.Status(204)
		} else {
			proxy(c)
		}
	})
}

func main() {
	startup()
	defer shutdown()
	r := gin.Default()
	store := sessions.NewCookieStore([]byte(env.SecretToken))
	r.Use(redirectHTTPS)
	r.Use(sessions.Sessions("labdb", store))
	r.POST("/api/verify", func(c *gin.Context) {
		email := auth.GetVerifiedIdentity(c.Query("token"))
		if email == "" {
			c.String(403, "Forbidden")
		} else {
			session := sessions.Default(c)
			session.Set("userID", email)
			session.Save()
			c.Redirect(303, "/")
		}
	})
	r.GET("/", proxy)

	// Below here, all routes require authorization.
	r.Use(requireAuthorization)

	r.GET("/search", func(c *gin.Context) {
		term := c.Query("term")
		seq := c.Query("seq")
		person := c.Query("person")
		includeSeq := false
		if seq == "1" {
			includeSeq = true
		}
		types := []string{}
		err := json.Unmarshal([]byte(c.Query("types")), &types)
		if err != nil || term == "" {
			c.String(400, "Invalid search query")
			return
		}
		results, err := search.Search(term, includeSeq, person, types)
		if err != nil {
			c.String(400, "Invalid search query")
			return
		}
		query := [][]interface{}{}
		for _, entity := range results {
			query = append(query, []interface{}{reflect.Indirect(reflect.ValueOf(entity)).Type().Name(), entity.GetID()})
		}
		// In addition to the searched items, we also send up the raw term to
		// check if it's a direct reference to an item (which we can't do in
		// the frontend yet).
		queryBytes, err := json.Marshal(map[string]interface{}{"items": query, "term": term})
		if err != nil {
			panic(err)
		}
		url, err := url.Parse(backendHost(c))
		if err != nil {
			panic(err)
		}
		req, err := http.NewRequest("POST", "/search_result", bytes.NewReader(queryBytes))
		if err != nil {
			panic(err)
		}
		setupBackendRequest(req, auth.CurrentUserID(c))
		req.Host = url.Host
		req.URL.Host = url.Host
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-CSRF-Token", c.Request.Header.Get("X-CSRF-Token"))
		for _, c := range c.Request.Cookies() {
			req.AddCookie(c)
		}
		if env.Dev {
			req.URL.Scheme = "http"
		} else {
			req.URL.Scheme = "https"
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Writer.Write(body)
	})

	routes.InstallAll(r)

	r.Use(proxy)

	if env.Dev {
		port := os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}
		r.Run(":" + port)
	} else {
		r.Run(":" + os.Getenv("PORT"))
	}
}
