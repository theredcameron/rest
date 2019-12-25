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

type dataWrapper struct {
	Data interface{} `json:"data"`
}

type errorWrapper struct {
	Error string `json:"error"`
}

type Router struct {
	muxRouter *mux.Router
}

type CookieMeta struct {
	File      string
	Path      string
	MaxAge    int
	Key       []byte
	StoreName string
}

func (this *Router) Start(port string) error {
	return http.ListenAndServe(":"+port, this.muxRouter)
}

var store *sqlitestore.SqliteStore

//func NewRouter(endpoints Endpoints, meta *CookieMeta) (*Router, error) {
func NewRouter(meta *CookieMeta, endpoints ...Endpoint) (*Router, error) {
	router := mux.NewRouter().StrictSlash(true)
	for _, endpoint := range endpoints {
		var handler http.Handler
		handler = NewHandlerFunc(endpoint.F, meta)
		router.
			Methods(endpoint.Method).
			Path(endpoint.Path).
			Name(endpoint.Description).
			Handler(handler)
	}
	if meta != nil {
		var err error
		store, err = sqlitestore.NewSqliteStore(meta.File, "sessions", meta.Path, meta.MaxAge, meta.Key)
		if err != nil {
			return nil, err
		}
	}
	return &Router{
		muxRouter: router,
	}, nil
}

func NewHandlerFunc(f func(*Request) (interface{}, error), meta *CookieMeta) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := NewRequest(r, meta)
		if err != nil {
			errorReturn(err, w)
			return
		}
		result, err := f(req)
		if err != nil {
			errorReturn(err, w)
			return
		}
		data := &dataWrapper{
			Data: result,
		}
		var session *sessions.Session
		if store != nil {
			session, err = store.Get(r, meta.StoreName)
			if err != nil {
				errorReturn(err, w)
				return
			}
			session.Values = req.getAllCookieValues()
			err = store.Save(r, w, session)
			if err != nil {
				errorReturn(err, w)
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

func errorReturn(err error, w http.ResponseWriter) {
	error_string := &errorWrapper{
		Error: fmt.Sprintf("Error: %v", err),
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(error_string); err != nil {
		panic(err)
	}
}
