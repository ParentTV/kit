package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ParentTV/kit/encoding"
	"github.com/ParentTV/kit/event_bus"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type User struct {
	Id          string   `json:"id"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Permissions []string `json:"permissions"`
	LicenseKey  string   `json:"license_key"`
}

type Auth struct {
	eb         *event_bus.EventBus
	url        string
	authString string
	token      string
}

func NewAuth(eb *event_bus.EventBus) *Auth {
	a := base64.StdEncoding.EncodeToString([]byte(os.Getenv("AUTH_SERVICE_USER") + ":" + os.Getenv("AUTH_SERVICE_PASS")))
	return &Auth{eb: eb, url: os.Getenv("AUTH_SERVICE_URL"), authString: "Basic " + a}
}

func (a *Auth) Login(user, pass string) (token string) {
	var b []byte
	as := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
	resp := a.get(a.url, as)
	defer resp.Body.Close()
	b, _ = ioutil.ReadAll(resp.Body)
	a.token = string(b)
	return string(b)
}

func (a *Auth) IdFromToken(token string) string {
	return "test"
}

func (a *Auth) IsAuthorized(id string, perm string) bool {
	if u, ok := a.getUser(id); ok {
		for _, v := range u.Permissions {
			if perm == v {
				return true
			}
		}
	}
	return false
}

func (a *Auth) HasLicense(id string, license string) bool {
	if u, ok := a.getUser(id); ok {
		return u.LicenseKey == license
	}
	return false
}

func (a *Auth) getUser(id string) (User, bool) {
	var u User
	url := fmt.Sprintf("%s/users/%s", a.url, id)
	r := a.get(url, "Bearer "+a.UseToken())
	if r == nil {
		return u, false
	}
	defer r.Body.Close()
	b, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(b, &u)
	return u, err != nil
}
func (a *Auth) UseToken() string {
	if a.token == "" {
		a.token = a.GetToken()
	}
	auth := encoding.ParseTokenString(a.token)
	if auth.Expiry.Before(time.Now()) {
		a.token = a.GetToken()
	}
	return a.token
}

func (a *Auth) GetToken() string {
	authlogin := a.url + "/auth"
	r := a.get(authlogin, a.authString)
	if r == nil {
		return ""
	}
	defer r.Body.Close()
	b, _ := ioutil.ReadAll(r.Body)
	return string(b)
}

func (a *Auth) get(url, authString string) *http.Response {
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
