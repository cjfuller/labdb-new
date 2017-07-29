package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"labdb.org/labdb/env"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var devProxyTarget = "http://localhost:3001"
var proxySuffix = "-backend.labdb.io"
var appID = "146923434465-alq7iagpanjvoag20smuirj0ivdtfldk.apps.googleusercontent.com"

type authResponse struct {
	Aud           string `json:"aud"`
	EmailVerified string `json:"email_verified"`
	Email         string `json:"email"`
}

func getVerifiedIdentity(token string) string {
	params := url.Values{}
	params.Set("id_token", token)
	url := url.URL{
		Scheme:   "https",
		Host:     "www.googleapis.com",
		Path:     "/oauth2/v3/tokeninfo",
		RawQuery: params.Encode(),
	}
	fmt.Println(url.String())
	resp, err := http.Post(url.String(), "text/plain", strings.NewReader(""))
	if err != nil || resp.StatusCode != 200 {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var authResp authResponse
	err = json.Unmarshal(body, &authResp)
	if err != nil {
		panic(err)
	}
	if strings.Contains(authResp.Aud, appID) && authResp.EmailVerified == "true" {
		return authResp.Email
	}
	return ""
}

func addAuthHeaders(userID string, h http.Header) {
	ts := time.Now().UTC().Format("2006-01-02T15:04:05")
	mac := hmac.New(sha256.New, []byte(env.SigningKey))
	mac.Write([]byte(userID + ts))
	result := hex.EncodeToString(mac.Sum(nil))
	h.Add("X-LabDB-UserId", userID)
	h.Add("X-LabDB-Signature", result)
	h.Add("X-LabDB-Signature-Timestamp", ts)
}

func proxy(c *gin.Context) {
	proxyTarget := devProxyTarget
	if env.Prod {
		proxyTarget = "https://" + strings.Replace(c.Request.Host, "labdb.io", proxySuffix, -1)
	}
	url, err := url.Parse(proxyTarget)
	if err != nil {
		panic(err)
	}

	for k := range c.Request.Header {
		if strings.HasPrefix(strings.ToLower(k), "cf-") {
			c.Request.Header.Del(k)
		}
	}
	c.Request.Header.Del("X-Forward-For")
	session := sessions.Default(c)
	maybeID := session.Get("userID")
	if maybeID != nil {
		addAuthHeaders(maybeID.(string), c.Request.Header)
	}

	p := httputil.NewSingleHostReverseProxy(url)
	p.ServeHTTP(c.Writer, c.Request)
}

func main() {
	env.Init()
	if env.Prod {
		gin.SetMode(gin.ReleaseMode)
	}
	// models.Init()
	// defer models.Shutdown()
	r := gin.Default()
	store := sessions.NewCookieStore([]byte(env.SecretToken))
	r.Use(sessions.Sessions("labdb", store))
	// r.GET("/plasmids/:id/next", func(c *gin.Context) {
	// 	id := c.Param("id")
	// 	pl := models.Plasmid{}
	// 	models.Db("dev").Where("id > ?", id).Order("id asc").First(&pl)
	// 	redirectID := id
	// 	if pl.ID != 0 {
	// 		redirectID = strconv.FormatUint(uint64(pl.ID), 10)
	// 	}
	// 	c.Redirect(307, fmt.Sprintf("/plasmids/%s", redirectID))
	// })
	// r.GET("/plasmids/:id/previous", func(c *gin.Context) {
	// 	id := c.Param("id")
	// 	pl := models.Plasmid{}
	// 	models.Db("dev").Where("id < ?", id).Order("id desc").First(&pl)
	// 	redirectID := id
	// 	if pl.ID != 0 {
	// 		redirectID = strconv.FormatUint(uint64(pl.ID), 10)
	// 	}
	// 	c.Redirect(307, fmt.Sprintf("/plasmids/%s", redirectID))
	// })
	r.POST("/api/verify", func(c *gin.Context) {
		email := getVerifiedIdentity(c.Query("token"))
		if email == "" {
			c.String(403, "Forbidden")
		} else {
			session := sessions.Default(c)
			session.Set("userID", email)
			session.Save()
			c.Redirect(303, "/")
		}
	})

	r.Use(proxy)

	if env.Dev {
		r.Run(":3000")
	} else {
		r.Run(":" + os.Getenv("PORT"))
	}
}
