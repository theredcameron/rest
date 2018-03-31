package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Description string
	Method      string
	Path        string
	f           func() interface{}
}

type Routes []Route

type Router mux.Router

func NewRouter(routes Routes) *Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
}
