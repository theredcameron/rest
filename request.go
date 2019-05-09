package rest

import (
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

func (this *Request) GetCookieValue(key interface{}) interface{} {
	if value, ok := this.cookieValues[key]; ok {
		return value
	}
	return nil
}

func (this *Request) getAllCookieValues() cookieValues {
	return this.cookieValues
}

func (this *Request) SetCookieValue(key, value interface{}) {
	if value == nil {
		delete(this.cookieValues, key)
		return
	}
	this.cookieValues[key] = value
}

func NewRequest(r *http.Request, meta *CookieMeta) (*Request, error) {
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

	var cookieVals cookieValues

	if meta != nil {
		session, err := store.Get(r, meta.StoreName)
		if err != nil {
			return nil, err
		}
		cookieVals = session.Values
	}

	return &Request{
		Vars:         mux.Vars(r),
		Body:         body,
		Params:       queries,
		cookieValues: cookieVals,
	}, nil
}
