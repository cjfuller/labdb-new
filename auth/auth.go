package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"labdb.org/labdb/env"
	"labdb.org/labdb/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var appID = "146923434465-alq7iagpanjvoag20smuirj0ivdtfldk.apps.googleusercontent.com"

type authResponse struct {
	Aud           string `json:"aud"`
	EmailVerified string `json:"email_verified"`
	Email         string `json:"email"`
}

func AddAuthHeaders(userID string, h http.Header) {
	ts := time.Now().UTC().Format("2006-01-02T15:04:05")
	mac := hmac.New(sha256.New, []byte(env.SigningKey))
	mac.Write([]byte(userID + ts))
	result := hex.EncodeToString(mac.Sum(nil))
	// TODO(colin): add this as a field to the normal log line?
	log.Printf("Verified user is: %s\n", userID)
	h.Add("X-LabDB-UserId", userID)
	h.Add("X-LabDB-Signature", result)
	h.Add("X-LabDB-Signature-Timestamp", ts)
}

func GetVerifiedIdentity(token string) string {
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
		fmt.Printf("Got error %+v, %+v\n", err, resp.StatusCode)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Body: %+v\n", string(body))

		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var authResp authResponse
	if env.Dev {
		fmt.Printf("Body: %+v\n", string(body))
	}
	err = json.Unmarshal(body, &authResp)
	if env.Dev {
		fmt.Printf("Auth resp: %+v\n", authResp)
	}
	if err != nil {
		panic(err)
	}
	if strings.Contains(authResp.Aud, appID) && authResp.EmailVerified == "true" {
		return authResp.Email
	}
	return ""
}

func CurrentUserID(c *gin.Context) string {
	session := sessions.Default(c)
	maybeID := session.Get("userID")
	if maybeID != nil {
		return maybeID.(string)
	}
	return ""
}

func CurrentUser(c *gin.Context) models.User {
	uid := CurrentUserID(c)
	return models.UserByEmail(uid)
}
