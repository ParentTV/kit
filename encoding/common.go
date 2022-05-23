package encoding

import (
	"errors"
	"github.com/golang/gddo/httputil/header"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
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

type StandardRequest struct {
	CorrelationId string
	Auth          Auth
	Params        map[string]string
	Query         url.Values
	Data          []byte
}

func NewStandardRequest(r *http.Request) StandardRequest {
	b, _, code := ParseJsonBody(r)
	if code == http.StatusBadRequest {
		return StandardRequest{}
	}
	p, q := ExtractParams(r)
	return StandardRequest{
		CorrelationId: ExtractCorrelationId(r),
		Auth:          ParseAuthHeader(r),
		Data:          b,
		Params:        p,
		Query:         q,
	}

}

func ExtractCorrelationId(r *http.Request) string {
	cid := r.Header.Get("x-correlation-id")
	if cid == "" {
		u, _ := uuid.NewUUID()
		cid = u.String()
	}
	return cid
}

func ExtractParams(r *http.Request) (map[string]string, url.Values) {
	return mux.Vars(r), r.URL.Query()
}
