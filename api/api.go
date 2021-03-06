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

	Response struct {
		Result int `json:"result"`
	}

	Welcome struct {
		Message string `json:"message"`
	}

	Api interface {
		Start()
	}
)

func New() Api {
	return &app{}
}

func (a *app) Start() {
	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler{})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&Welcome{Message: "Welcome to the home page!"})
	})

	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func (apiHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&Response{Result: 1})
}
