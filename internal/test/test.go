package test

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
}

// TestWithoutAnswers is the same as Test but without the answers.
// Used to send the test to the client.
type TestWithoutAnswers struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	ShortDescription       string `json:"shortDescription"`
	Description            string `json:"description"`
	QuestionMetadata       `json:"metadata"`
	Instructions           `json:"instructions"`
	QuestionsWithoutAnswer []QuestionWithoutAnswer `json:"questions"`
}

type QuestionWithoutAnswer struct {
	ID      string   `json:"id"`
	Text    string   `json:"text"`
	Options []Option `json:"options"`
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
	ID              string   `json:"id"`
	Text            string   `json:"text"`
	Options         []Option `json:"options"`
	CorrectOptionId string   `json:"correctOptionId"`
	Explanation     string   `json:"explanation"`
}

// ID and text of an option.
type Option struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

func (t Test) removeAnswers() TestWithoutAnswers {
	var tmp TestWithoutAnswers
	tmp.ID = t.ID
	tmp.Name = t.Name
	tmp.ShortDescription = t.ShortDescription
	tmp.Description = t.Description
	tmp.QuestionMetadata = t.QuestionMetadata
	tmp.Instructions = t.Instructions
	tmp.QuestionsWithoutAnswer = []QuestionWithoutAnswer{}
	for _, q := range t.Questions {
		questionWithoutAnswer := QuestionWithoutAnswer{
			ID:      q.ID,
			Text:    q.Text,
			Options: q.Options,
		}
		tmp.QuestionsWithoutAnswer = append(tmp.QuestionsWithoutAnswer, questionWithoutAnswer)
	}

	return tmp
}

// removeAnswersFromTests is a helper function that removes answers from a slice of Tests.
func removeAnswersFromTests(tests []Test) []TestWithoutAnswers {
	var tmp []TestWithoutAnswers
	for _, t := range tests {
		tmp = append(tmp, t.removeAnswers())
	}
	return tmp
}
