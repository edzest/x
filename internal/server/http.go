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
	r.Use(auth)
	r.HandleFunc("/tests", th.ListTests).Methods("GET")
	r.HandleFunc("/tests", th.CreateTest).Methods("POST")
	r.HandleFunc("/tests/{id}", th.GetTest).Methods("GET")
	r.HandleFunc("/tests/{id}:evaluate", th.EvaluateTest).Methods("POST")

	r.HandleFunc("/health", handleHealthCheck)

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintln(w, `{"status":"ok"}`)
}
