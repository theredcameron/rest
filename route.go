package rest

import (
	"encoding/json"
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
	Error error `json:"error"`
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
		req := NewRequest(r)
		result, err := f(req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusInternalServerError)
		}
		data := &DataWrapper{
			Data: result,
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(data); err != nil {
			error_string := &ErrorWrapper{
				Error: err,
			}
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(error_string); err != nil {
				panic(err)
			}
		}
	}
}
