package server

import (
	"encoding/json"
	"fmt"
)

// Test is the primary object of this server and contains the id, name
// and description of test along with the list of questions and their
// correct answers.
type Test struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
	QuestionMetadata `json:"metadata"`
	Instructions     `json:"instructions"`
	Questions        []Question `json:"questions"`
	Answers          []Answer   `json:"answers"`
}

// QuestionMetadata contains the summary of the test.
type QuestionMetadata struct {
	TotalQuestions int8  `json:"total_questions"`
	TotalMarks     int8  `json:"total_marks"`
	TotalTime      int16 `json:"total_time"`
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

// Answer contains id of question and its correct option.
type Answer struct {
	QId string `json:"q_id"`
	AId string `json:"a_id"`
}

// Custom func to omit Answers while marshalling a Test
func (t Test) MarshalJSON() ([]byte, error) {
	var tmp struct {
		ID               string `json:"id"`
		Name             string `json:"name"`
		ShortDescription string `json:"short_description"`
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

// Temporary in memory db.
type DB struct {
	records map[string]Test
}

func NewDB() *DB {
	return &DB{
		make(map[string]Test),
	}
}

var ErrorTestNotFound = fmt.Errorf("test not found")

func (db *DB) insert(t Test) {
	db.records[t.ID] = t
}

func (db *DB) get(id string) (Test, error) {
	if t, ok := db.records[id]; ok {
		return t, nil
	}
	return Test{}, ErrorTestNotFound
}
