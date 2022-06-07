package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ParentTV/kit/event_bus"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type User struct {
	Id          string   `json:"id"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Permissions []string `json:"permissions"`
}

type Auth struct {
	eb         *event_bus.EventBus
	url        string
	authString string
}

func NewAuth(eb *event_bus.EventBus) *Auth {
	a := base64.StdEncoding.EncodeToString([]byte(os.Getenv("AUTH_SERVICE_USER") + ":" + os.Getenv("AUTH_SERVICE_PASS")))
	return &Auth{eb, os.Getenv("AUTH_SERVICE_URL"), "Basic " + a}
}

func (a *Auth) Login(user, pass string) (token string) {
	var b []byte
	as := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
	resp := a.login(a.url, as)
	defer resp.Body.Close()
	b, _ = ioutil.ReadAll(resp.Body)
	return string(b)
}

func (a *Auth) IdFromToken(token string) string {
	return "test"
}

func (a *Auth) IsAuthorized(id string, perm string) bool {
	url := fmt.Sprintf("%s/users/%s", a.url, id)
	r := a.login(url, a.authString)
	if r == nil {
		return false
	}
	defer r.Body.Close()
	b, _ := ioutil.ReadAll(r.Body)
	var u User
	err := json.Unmarshal(b, &u)
	if err != nil {
		return false
	}
	for _, v := range u.Permissions {
		if perm == v {
			return true
		}
	}
	return false
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
		return nil
	}
	return resp
}

func (a *Auth) RegisterPermissions(service string, perms []string) {
	e := a.eb.NewEvent()
	e.SetTopic(service + ".permissions.register")
	e.SetData(perms)
	e.SetOrigin(service)
	a.eb.Publish(e)
}
