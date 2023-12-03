package test

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

// Answer contains id of question and its correct option.
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

// TestStore is an interface that provides reading and writing
// a test in a database.
type TestStore interface {
	insert(t Test)
	get(id string) (Test, error)
}

// Temporary in memory db satisfies TestStore interface.
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
