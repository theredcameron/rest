package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Endpoint struct {
	Description string
	Method      string
	Path        string
	F           func(*Request) (interface{}, error)
}

type Endpoints []Endpoint

type DataWrapper struct {
	Data interface{} `json:"data"`
}

type ErrorWrapper struct {
	Error string `json:"error"`
}

type Router struct {
	muxRouter *mux.Router
}

func (this *Router) Start(port string) error {
	return http.ListenAndServe(":"+port, this.muxRouter)
}

func NewRouter(endpoints Endpoints) *Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, endpoint := range endpoints {
		var handler http.Handler
		handler = NewHandlerFunc(endpoint.F)
		router.
			Methods(endpoint.Method).
			Path(endpoint.Path).
			Name(endpoint.Description).
			Handler(handler)
	}
	return &Router{
		muxRouter: router,
	}
}

func NewHandlerFunc(f func(*Request) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := NewRequest(r)
		if err != nil {
			ErrorReturn(err, w)
			return
		}
		result, err := f(req)
		if err != nil {
			ErrorReturn(err, w)
			return
		}
		data := &DataWrapper{
			Data: result,
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(data); err != nil {
			panic(err)
		}

	}
}

func ErrorReturn(err error, w http.ResponseWriter) {
	error_string := &ErrorWrapper{
		Error: fmt.Sprintf("Error: %v", err),
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(error_string); err != nil {
		panic(err)
	}
}
