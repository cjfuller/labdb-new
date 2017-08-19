package main

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"labdb.org/labdb/auth"
	"labdb.org/labdb/env"
	"labdb.org/labdb/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var devProxyTarget = "http://localhost:3001"
var proxySuffix = "-backend.labdb.io"

func proxy(c *gin.Context) {
	proxyTarget := devProxyTarget
	if env.Prod {
		proxyTarget = "https://" + strings.Replace(c.Request.Host, ".labdb.io", proxySuffix, -1)
	}
	url, err := url.Parse(proxyTarget)
	if err != nil {
		panic(err)
	}

	if c.Request.Header.Get("X-Labdb-Forwarded") == "true" {
		c.String(400, "Stuck in a recursive proxy loop.")
		return
	}

	for k := range c.Request.Header {
		if strings.HasPrefix(strings.ToLower(k), "cf-") {
			c.Request.Header.Del(k)
		}
	}
	c.Request.Header.Del("X-Forwarded-For")
	c.Request.Header.Add("X-Labdb-Forwarded", "true")
	maybeID := auth.CurrentUserID(c)
	if maybeID != "" {
		auth.AddAuthHeaders(maybeID, c.Request.Header)
	}
	// Unset the host, as the proxy sets the host using URL.Host, but the value on
	// the request itself will override it if present.
	c.Request.Host = ""
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
	env.Init()
	if env.Prod {
		gin.SetMode(gin.ReleaseMode)
	}
	models.Init()
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

	r.GET("/:model/:id/next", func(c *gin.Context) {
		cls := c.Param("model")
		id := c.Param("id")
		redirectID := models.NextID(cls, id)
		c.Redirect(307, fmt.Sprintf("/%s/%s", cls, redirectID))
	})
	r.GET("/:model/:id/previous", func(c *gin.Context) {
		cls := c.Param("model")
		id := c.Param("id")
		redirectID := models.PrevID(cls, id)
		c.Redirect(307, fmt.Sprintf("/%s/%s", cls, redirectID))
	})

	r.Use(proxy)

	if env.Dev {
		r.Run(":3000")
	} else {
		r.Run(":" + os.Getenv("PORT"))
	}
}
