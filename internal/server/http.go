package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func NewHTTPServer(addr string) *http.Server {
	httpsrv := newHTTPServer()
	r := mux.NewRouter()
	r.HandleFunc("/tests", httpsrv.handleCreateTest).Methods("POST")
	r.HandleFunc("/tests/{id}", httpsrv.handleGetTest).Methods("GET")
	r.HandleFunc("/tests/{id}:evaluate", httpsrv.handleEvaluateTest).Methods("POST")

	r.HandleFunc("/health", handleHealthCheck)

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

type httpServer struct {
	Db          *DB
	evalService *EvaluationService
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Db: NewDB(),
	}
}

func (s *httpServer) handleGetTest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	t, err := s.Db.get(id)
	if err == ErrorTestNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) handleCreateTest(w http.ResponseWriter, r *http.Request) {
	var t Test
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.Db.insert(t)
	w.WriteHeader(http.StatusOK)
}

func (s *httpServer) handleEvaluateTest(w http.ResponseWriter, r *http.Request) {
	var t TestEvalRequest

	vars := mux.Vars(r)
	id := vars["id"]
	t.TestId = id

	var answers []Answer
	err := json.NewDecoder(r.Body).Decode(&answers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Answers = answers

	res := s.evalService.evaluate(t)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintln(w, `{"status":"ok"}`)
}
