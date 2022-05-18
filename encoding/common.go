package encoding

import (
	"errors"
	"github.com/golang/gddo/httputil/header"
	"net/http"
)

func headerCheck(t string, r *http.Request) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != t {
			return errors.New("Content-Type header is not " + t)
		}
	}
	return nil
}
