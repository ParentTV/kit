package encoding

import "net/http"

type Auth struct {
	User     string
	Password string
}

func ParseAuthHeader(r *http.Request) Auth {
	u, p, ok := r.BasicAuth()
	if !ok {
		return Auth{}
	}
	return Auth{u, p}
}
