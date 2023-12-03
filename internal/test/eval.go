package test

import "log"

type EvaluationService struct {
	testStore TestStore
}

func NewEvaluationService() *EvaluationService {
	return &EvaluationService{
		testStore: NewTempTestDB(),
	}
}

type Result struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
	Total int    `json:"total"`
}

type TestEvalRequest struct {
	Name    string   `json:"name"`
	TestId  string   `json:"testId"`
	Answers []Answer `json:"answers"`
}

func (e *EvaluationService) evaluate(r TestEvalRequest) (*Result, error) {
	log.Printf("Fetching test with id %s\n", r.TestId)
	t, err := e.testStore.get(r.TestId)
	if err != nil {
		log.Printf("Test with id %s not found", r.TestId)
		log.Println(err.Error())
		return nil, ErrorTestNotFound
	}
	score := 0
	for _, ans := range r.Answers {
		if aId, ok := t.getAnswer(ans.QId); ok {
			if ans.AId == aId {
				score++
			}
		}
	}
	return &Result{
		Name:  t.Name,
		Score: score,
		Total: len(t.Answers),
	}, nil
}
