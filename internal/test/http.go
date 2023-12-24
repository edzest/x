package test

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type testHandler struct {
	testStore   TestStore
	evalService *EvaluationService
}

func NewHttpHandler() (*testHandler, error) {
	client, err := initializeFirebaseDb()
	if err != nil {
		log.Fatal("error initializing firebase db client:", err)
		return nil, err
	}
	testStore := NewFbTestDB(client)
	return &testHandler{
		testStore:   testStore,
		evalService: NewEvaluationService(testStore),
	}, nil
}

// ListTests returns a list of all tests.
// todo: Add pagination params.
func (h *testHandler) ListTests(w http.ResponseWriter, r *http.Request) {
	tests, err := h.testStore.list()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ListTestsResponse{
		TestsWithoutAnswers: removeAnswersFromTests(tests),
		Total:               len(tests),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	err = json.NewEncoder(w).Encode(t.removeAnswers())
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

// ListTestsResponse is the response returned by ListTests handler.
type ListTestsResponse struct {
	TestsWithoutAnswers []TestWithoutAnswers `json:"tests"`
	Total               int                  `json:"total"`
}
