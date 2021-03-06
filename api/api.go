package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type (
	app struct{}

	apiHandler struct{}

	rootHandler struct{}

	Response struct {
		Result int `json:"result"`
	}

	Welcome struct {
		Message string `json:"message"`
	}

	Api interface {
		Start()
	}

	JSONHandler interface {
		ServeHTTP(http.ResponseWriter, *http.Request)
		Result(w http.ResponseWriter, r *http.Request) interface{}
	}
)

func New() Api {
	return &app{}
}

func (a *app) Start() {
	mux := http.NewServeMux()
	mux.Handle("/api/", jsonMiddleware(apiHandler{}))
	mux.Handle("/", jsonMiddleware(rootHandler{}))

	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func (apiHandler) Result(http.ResponseWriter, *http.Request) interface{} {
	return &Response{Result: 1}
}

func (rootHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func (rootHandler) Result(w http.ResponseWriter, r *http.Request) interface{} {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return nil
	}
	return &Welcome{Message: "Welcome to the home page!"}
}

func jsonMiddleware(next JSONHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result := next.Result(w, r)
		if result != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		}
	})
}
