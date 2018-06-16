package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/michaeljs1990/sqlitestore"
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

type Auth struct {
	File   string
	Table  string
	Path   string
	MaxAge int
	Key    []byte
}

func (this *Router) Start(port string) error {
	return http.ListenAndServe(":"+port, this.muxRouter)
}

var store *sqlitestore.SqliteStore

func NewRouter(endpoints Endpoints, auth *Auth) (*Router, error) {
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
	if auth != nil {
		var err error
		store, err = sqlitestore.NewSqliteStore(auth.File, "sessions", auth.Path, auth.MaxAge, auth.Key)
		if err != nil {
			return nil, err
		}
	}
	return &Router{
		muxRouter: router,
	}, nil
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
		var session *sessions.Session
		if store != nil {
			session, err = store.Get(r, "authentication")
			if err != nil {
				ErrorReturn(err, w)
				return
			}
			session.Values = req.CookieValues
			err = session.Save(r, w)
			if err != nil {
				ErrorReturn(err, w)
				return
			}
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
