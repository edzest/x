package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func NewHTTPServer(addr string) *http.Server {
	httpsrv := newHTTPServer()
	r := mux.NewRouter()
	r.HandleFunc("/", httpsrv.handleCreateTest).Methods("POST")
	r.HandleFunc("/", httpsrv.handleGetTest).Methods("GET")

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

type httpServer struct {
	Db *DB
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Db: NewDB(),
	}
}

func (s *httpServer) handleGetTest(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	ids := params["id"]
	t, err := s.Db.get(ids[0])
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
