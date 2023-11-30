package server

type EvaluationService struct{}

type Result struct {
	Grade string `json:"grade"`
}

type TestEvalRequest struct {
	TestId  string   `json:"testId"`
	Answers []Answer `json:"answers"`
}

func (*EvaluationService) evaluate(t TestEvalRequest) *Result {
	return &Result{
		Grade: "A",
	}
}
