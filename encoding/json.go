package encoding

import (
	"github.com/gorilla/mux"
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
	err := headerCheck("application/json", r)
	if err != nil {
		return JsonData{Error: err, ErrorCode: http.StatusUnsupportedMediaType}
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return JsonData{Error: err, ErrorCode: http.StatusBadRequest}
	}
	jsonData := JsonData{
		Params: mux.Vars(r),
		Query:  r.URL.Query(),
		Data:   b,
	}
	return jsonData
}
