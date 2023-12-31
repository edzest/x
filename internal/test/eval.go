package test

import "log"

type EvaluationService struct {
	testStore TestStore
}

func NewEvaluationService(testStore TestStore) *EvaluationService {
	return &EvaluationService{
		testStore: testStore,
	}
}

type Result struct {
	Name            string           `json:"name"`
	Score           int              `json:"score"`
	Total           int              `json:"total"`
	GuidedSolutions []GuidedSolution `json:"answerSheet"`
}

// GuidedSolution contains the question, correct answer with explanation and the selected answer.
type GuidedSolution struct {
	Question         Question `json:"question"`
	SelectedAnswerId string   `json:"selectedAnswerId"`
}

type TestEvalRequest struct {
	Name            string            `json:"name"`
	TestId          string            `json:"testId"`
	SelectedAnswers map[string]string `json:"answers"`
}

func (e *EvaluationService) evaluate(r TestEvalRequest) (*Result, error) {
	log.Printf("Fetching test with id %s for evaluation\n", r.TestId)
	t, err := e.testStore.get(r.TestId)
	if err != nil {
		log.Printf("Test with id %s not found", r.TestId)
		log.Println(err.Error())
		return nil, ErrorTestNotFound
	}
	score := 0
	guidedSolutions := make([]GuidedSolution, 0)

	for _, q := range t.Questions {
		aId, attempted := r.SelectedAnswers[q.ID]

		if attempted && aId == q.CorrectOptionId {
			score++
		}

		guidedSolution := GuidedSolution{
			Question:         q,
			SelectedAnswerId: aId,
		}
		guidedSolutions = append(guidedSolutions, guidedSolution)
	}

	return &Result{
		Name:            t.Name,
		Score:           score,
		Total:           len(t.Questions),
		GuidedSolutions: guidedSolutions,
	}, nil
}
