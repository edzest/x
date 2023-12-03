package test

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type testHandler struct {
	testStore   TestStore
	evalService *EvaluationService
}

func NewHttpHandler() *testHandler {
	return &testHandler{
		testStore:   NewTempTestDB(),
		evalService: NewEvaluationService(),
	}
}

func (h *testHandler) GetTest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	t, err := h.testStore.get(id)
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

func (h *testHandler) CreateTest(w http.ResponseWriter, r *http.Request) {
	var t Test
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.testStore.insert(t)
	w.WriteHeader(http.StatusOK)
}

func (h *testHandler) EvaluateTest(w http.ResponseWriter, r *http.Request) {
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

	res, err := h.evalService.evaluate(t)

	if err == ErrorTestNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
