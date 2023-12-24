package test

import "encoding/json"

// Test is the primary object of this server and contains the id, name
// and description of test along with the list of questions and their
// correct answers.
type Test struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ShortDescription string `json:"shortDescription"`
	Description      string `json:"description"`
	QuestionMetadata `json:"metadata"`
	Instructions     `json:"instructions"`
	Questions        []Question `json:"questions"`
	Answers          []Answer   `json:"answers"`
}

// QuestionMetadata contains the summary of the test.
type QuestionMetadata struct {
	TotalQuestions int8  `json:"totalQuestions"`
	TotalMarks     int8  `json:"totalMarks"`
	TotalTime      int16 `json:"totalTime"`
}

// Test instructions stored as rich text html.
type Instructions struct {
	Body string `json:"body"`
}

// Question contains text of the question and its answers.
type Question struct {
	ID      string   `json:"id"`
	Text    string   `json:"text"`
	Options []Option `json:"options"`
}

// ID and text of an option.
type Option struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// Answer contains ids of question and its correct option.
type Answer struct {
	QId string `json:"qId"`
	AId string `json:"aId"`
}

// Custom func to omit Answers while marshalling a Test
func (t Test) MarshalJSON() ([]byte, error) {
	var tmp struct {
		ID               string `json:"id"`
		Name             string `json:"name"`
		ShortDescription string `json:"shortDescription"`
		Description      string `json:"description"`
		QuestionMetadata `json:"metadata"`
		Instructions     `json:"instructions"`
		Questions        []Question `json:"questions"`
	}
	tmp.ID = t.ID
	tmp.Name = t.Name
	tmp.ShortDescription = t.ShortDescription
	tmp.Description = t.Description
	tmp.QuestionMetadata = t.QuestionMetadata
	tmp.Instructions = t.Instructions
	tmp.Questions = t.Questions
	return json.Marshal(&tmp)
}

func (t Test) getAnswer(qid string) (string, bool) {
	for _, q := range t.Answers {
		if q.QId == qid {
			return q.AId, true
		}
	}
	return "", false
}
