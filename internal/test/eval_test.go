package test

import "testing"

func TestSimpleEvaluation(t *testing.T) {
	store := make(map[string]Test)
	evalService := &EvaluationService{
		testStore: &TempTestDB{store},
	}
	simpleTest := Test{
		ID:               "1234",
		Name:             "Simple Test",
		ShortDescription: "",
		Description:      "",
		QuestionMetadata: QuestionMetadata{
			TotalQuestions: 2,
			TotalMarks:     2,
			TotalTime:      120,
		},
		Instructions: Instructions{
			Body: "",
		},
		Questions: []Question{
			{
				ID:   "1",
				Text: "What is the atomic number of Hydrogen?",
				Options: []Option{
					{ID: "1", Text: "1"},
					{ID: "2", Text: "2"},
					{ID: "3", Text: "10"},
					{ID: "4", Text: "12"},
				},
			},
			{
				ID:   "2",
				Text: "Who discovered Radium?",
				Options: []Option{
					{ID: "1", Text: "Madam Curie"},
					{ID: "2", Text: "C V Raman"},
					{ID: "3", Text: "Michael Faraday"},
					{ID: "4", Text: "James Maxwell"},
				},
			},
		},
		Answers: []Answer{
			{QId: "1", AId: "1"},
			{QId: "2", AId: "1"},
		},
	}
	store["1234"] = simpleTest

	request := TestEvalRequest{
		Name:   simpleTest.Name,
		TestId: simpleTest.ID,
		Answers: []Answer{
			{QId: "1", AId: "1"},
			{QId: "2", AId: "1"},
		},
	}

	res, err := evalService.evaluate(request)

	if err != nil {
		t.Error(err)
	}

	if res.Score != 2 {
		t.Errorf("expected total=%d, but got %d", 2, res.Score)
	}

	if res.Total != 2 {
		t.Errorf("expected score=%d, but got %d", 2, res.Total)
	}
}
