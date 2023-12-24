package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetReturnsTestWithoutAnswers(t *testing.T) {
	// create a test handler with a temp test store
	testStore := NewTempTestDB()
	th := &testHandler{
		testStore:   testStore,
		evalService: NewEvaluationService(testStore),
	}

	// add a test to the temp test store
	sampleTest := getSampleTest("1")
	th.testStore.insert(sampleTest)

	// create a request to get the test
	req, err := http.NewRequest("GET", "/tests/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// set mux vars
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(th.GetTest)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v\n", status, http.StatusOK)
	}

	// check the response body does not contain answers
	want := sampleTest.removeAnswers()

	var got TestWithoutAnswers
	json.Unmarshal(rr.Body.Bytes(), &got)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("handler returned unexpected body: got %v, want %v\n", got, want)
	}
}

func TestListTestsReturnsTestWithoutAnswers(t *testing.T) {
	// create a test handler with a temp test store
	testStore := NewTempTestDB()
	th := &testHandler{
		testStore:   testStore,
		evalService: NewEvaluationService(testStore),
	}

	// add two tests to the temp test store
	testOne := getSampleTest("1")
	testTwo := getSampleTest("2")
	th.testStore.insert(testOne)
	th.testStore.insert(testTwo)

	// create a request to get the test
	req, err := http.NewRequest("GET", "/tests", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(th.ListTests)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v\n", status, http.StatusOK)
	}

	// check the response body does not contain answers
	want := []TestWithoutAnswers{testOne.removeAnswers(), testTwo.removeAnswers()}

	var got ListTestsResponse
	json.Unmarshal(rr.Body.Bytes(), &got)

	if !reflect.DeepEqual(got.TestsWithoutAnswers, want) {
		t.Errorf("handler returned unexpected body: got %v, want %v\n", got, want)
	}
}

// getSampleTest returns a sample Test.
func getSampleTest(id string) Test {
	return Test{
		ID:               id,
		Name:             "Test 1",
		ShortDescription: "Test 1",
		Description:      "Test 1",
		QuestionMetadata: QuestionMetadata{
			TotalQuestions: 1,
			TotalMarks:     1,
			TotalTime:      10,
		},
		Instructions: Instructions{
			Body: "<p>Test 1</p>",
		},
		Questions: []Question{
			{
				ID:   "1",
				Text: "Question 1",
				Options: []Option{
					{
						ID:   "1",
						Text: "Option 1",
					},
					{
						ID:   "2",
						Text: "Option 2",
					},
				},
			},
		},
		Answers: []Answer{
			{
				QId: "1",
				AId: "1",
			},
		},
	}
}
