package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/ParentTV/kit/event_bus"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Auth struct {
	eb         *event_bus.EventBus
	url        string
	authString string
}

func NewAuth(eb *event_bus.EventBus) *Auth {
	a := base64.StdEncoding.EncodeToString([]byte(os.Getenv("AUTH_SERVICE_USER") + ":" + os.Getenv("AUTH_SERVICE_PASS")))
	return &Auth{eb, os.Getenv("AUTH_SERVICE_URL"), "Basic " + a}
}

func (a *Auth) Login(user, pass string) string {
	var b []byte
	as := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
	resp := a.login(a.url, as)
	defer resp.Body.Close()
	b, _ = ioutil.ReadAll(resp.Body)
	return string(b)

}

func (a *Auth) IsAuthorized(id string, perm string) bool {
	url := fmt.Sprintf("%s/%s?perm=%s", a.url, id, perm)
	return a.login(url, a.authString).StatusCode == 200
}

func (a *Auth) login(url, authString string) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error on request.\n[ERROR] -", err)
	}
	// add authorization header to the req
	req.Header.Add("Authorization", authString)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	return resp
}

func (a *Auth) RegisterPermissions(perms []string) {
	a.eb.Publish("permissions.register", perms)
}
