package rest

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type Request struct {
	Vars         map[string]string
	Body         []byte
	Params       map[string]string
	cookieValues cookieValues
}

type cookieValues map[interface{}]interface{}

func (this *Request) GetCookie(key interface{}) (interface{}, error) {
	if value, ok := this.cookieValues[key]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("entry not found")
}

func (this *Request) GetAllCookies() cookieValues {
	return this.cookieValues
}

func (this *Request) SetCookie(key, value interface{}) {
	this.cookieValues[key] = value
}

func NewRequest(r *http.Request) (*Request, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 9999999))
	if err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}
	queries := make(map[string]string)
	query_map := r.URL.Query()
	for index, value := range query_map {
		param, err := url.QueryUnescape(value[0])
		if err != nil {
			queries[index] = value[0]
			continue
		}
		queries[index] = param
	}

	session, err := store.Get(r, "authentication")
	if err != nil {
		return nil, err
	}

	var cookieVals cookieValues
	cookieVals = session.Values

	return &Request{
		Vars:         mux.Vars(r),
		Body:         body,
		Params:       queries,
		cookieValues: cookieVals,
	}, nil
}
