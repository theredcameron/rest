package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type InternalRoute struct {
	Description string
	Method      string
	Path        string
	f           http.HandlerFunc
}

type Route struct {
	Description string
	Method      string
	Path        string
	F           func(*Request) (interface{}, error)
}

type Routes []Route

type DataWrapper struct {
	Data interface{} `json:"data"`
}

type ErrorWrapper struct {
	Error string `json:"error"`
}

type InternalRoutes []InternalRoute

type Router struct {
	muxRouter *mux.Router
}

func (this *Router) Start(port string) error {
	return http.ListenAndServe(port, this.muxRouter)
}

func NewRouter(routes Routes) *Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = NewHandlerFunc(route.F)
		//handler = Logger(handler, route.Name) Use this method of 'logging' later for user authentication
		router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Description).
			Handler(handler)
	}
	return &Router{
		muxRouter: router,
	}
}

func NewHandlerFunc(f func(r *Request) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := NewRequest(r)
		result, err := f(req)
		if err != nil {
			ErrorReturn(err, w)
			return
		}
		data := &DataWrapper{
			Data: result,
		}

		if err := json.NewEncoder(w).Encode(data); err != nil {
			ErrorReturn(err, w)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
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
