package server

import (
	"fmt"
	"log"
	"net/http"

	"edzest.org/x/internal/test"

	"github.com/gorilla/mux"
)

func NewHTTPServer(addr string) *http.Server {
	th, err := test.NewHttpHandler()
	if err != nil {
		log.Fatal("error initializing http handler:", err)
		return nil
	}
	r := mux.NewRouter()

	r.HandleFunc("/health", handleHealthCheck)

	tr := r.PathPrefix("/tests").Subrouter()
	tr.HandleFunc("", th.ListTests).Methods("GET")
	tr.HandleFunc("", th.CreateTest).Methods("POST")
	tr.HandleFunc("/{id}", th.GetTest).Methods("GET")
	tr.HandleFunc("/{id}:evaluate", th.EvaluateTest).Methods("POST")
	tr.Use(auth)

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintln(w, `{"status":"ok"}`)
}
