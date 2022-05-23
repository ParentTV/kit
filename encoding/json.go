package encoding

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type JsonData struct {
	Params    map[string]string
	Query     url.Values
	Data      []byte
	Error     error
	ErrorCode int
}

func ParseJsonData(r *http.Request) JsonData {
	b, err, code := ParseJsonBody(r)
	if err != nil {
		return JsonData{Error: err, ErrorCode: code}
	}
	p, q := ExtractParams(r)
	jsonData := JsonData{
		Data:   b,
		Params: p,
		Query:  q,
	}
	return jsonData
}

func ParseJsonBody(r *http.Request) (b []byte, err error, code int) {
	defer r.Body.Close()
	err = headerCheck("application/json", r)
	if err != nil {
		return b, err, http.StatusUnsupportedMediaType
	}
	b, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return b, err, http.StatusBadRequest
	}
	return b, nil, http.StatusOK
}
