package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/ParentTV/kit/event_bus"
	"log"
	"net/http"
	"os"
)

type Auth struct {
	eb         *event_bus.EventBus
	url        string
	authstring string
}

func NewAuth(eb *event_bus.EventBus) *Auth {
	a := base64.StdEncoding.EncodeToString([]byte(os.Getenv("AUTH_SERVICE_USER") + ":" + os.Getenv("AUTH_SERVICE_PASS")))
	return &Auth{eb, os.Getenv("AUTH_SERVICE_URL"), "Basic " + a}
}

func (a *Auth) IsAuthorized(id string, perm string) bool {
	reqUrl := fmt.Sprintf("%s?id=%s&perm=%s", a.url, id, perm)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.Println("Error on request.\n[ERROR] -", err)
	}
	// add authorization header to the req
	req.Header.Add("Authorization", a.authstring)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	return resp.StatusCode == 200
}

func (a *Auth) RegisterPermissions(perms []string) {
	a.eb.Publish("permissions.register", perms)
}
