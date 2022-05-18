package encoding

import "net/http"

type Auth struct {
	User string
}

func ParseAuthHeader(r *http.Request) Auth {
	return Auth{}
}
