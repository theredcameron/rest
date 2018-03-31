package rest

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Request struct {
	Vars map[string]string
	Body []byte
}

func NewRequest(r *http.Request) (*Request, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 9999999))
	if err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}
	return &Request{
		Vars: mux.Vars(r),
		Body: body,
	}, nil
}
