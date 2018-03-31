package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Request struct {
	Vars map[string]string
}

func NewRequest(r *http.Request) *Request {
	return &Request{
		Vars: mux.Vars(r),
	}
}
