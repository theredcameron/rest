package rest

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type Request struct {
	Vars   map[string]string
	Body   []byte
	Params map[string]string
	Auth   *Authentication
}

type Authentication struct {
	Username string
	Password string
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

	var auth *Authentication

	username, password, ok := r.BasicAuth()
	if ok {
		auth = &Authentication{
			Username: username,
			Password: password,
		}
	}

	return &Request{
		Vars:   mux.Vars(r),
		Body:   body,
		Params: queries,
		Auth: auth,
	}, nil
}
