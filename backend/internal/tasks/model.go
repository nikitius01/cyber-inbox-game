package tasks

import "time"

type Difficulty string

const (
	Easy      Difficulty = "easy"
	Medium    Difficulty = "medium"
	Hard      Difficulty = "hard"
	Nightmare Difficulty = "nightmare"
	AI        Difficulty = "AI"
)

type Link struct {
	VisibleText  string `json:"visibleText"`
	ActualURL    string `json:"actualUrl"`
	Domain       string `json:"domain"`
	Protocol     string `json:"protocol"`
	IsShortened  bool   `json:"isShortened"`
	IsSuspicious bool   `json:"isSuspicious"`
	RiskReason   string `json:"riskReason,omitempty"`
}

type Attachment struct {
	FileName     string `json:"fileName"`
	DisplayName  string `json:"displayName"`
	Extension    string `json:"extension"`
	MimeType     string `json:"mimeType"`
	SizeKB       int    `json:"sizeKb"`
	Hash         string `json:"hash"`
	IsSuspicious bool   `json:"isSuspicious"`
	RiskReason   string `json:"riskReason,omitempty"`
}

type RawEmail struct {
	Headers  map[string]any `json:"headers"`
	BodyText string         `json:"bodyText"`
	BodyHTML string         `json:"bodyHtml"`
	Source   string         `json:"source"`
}

type RedFlag struct {
	Type        string `json:"type"`
	Field       string `json:"field"`
	Value       string `json:"value"`
	Explanation string `json:"explanation"`
}

type Task struct {
	ID          string       `json:"id"`
	Category    Difficulty   `json:"category"`
	SenderName  string       `json:"senderName"`
	SenderEmail string       `json:"senderEmail"`
	Subject     string       `json:"subject"`
	Body        string       `json:"body"`
	IsPhishing  bool         `json:"isPhishing"`
	Links       []Link       `json:"links"`
	Attachments []Attachment `json:"attachments"`
	Raw         RawEmail     `json:"raw"`
	RedFlags    []RedFlag    `json:"-"`
	Explanation string       `json:"-"`
	CreatedAt   time.Time    `json:"createdAt"`
}

type PublicTask struct {
	ID          string         `json:"id"`
	Category    Difficulty     `json:"category"`
	SenderName  string         `json:"senderName"`
	SenderEmail string         `json:"senderEmail"`
	Subject     string         `json:"subject"`
	Body        string         `json:"body"`
	Links       []Link         `json:"links"`
	Attachments []Attachment   `json:"attachments"`
	Raw         RawEmail       `json:"raw"`
	CreatedAt   time.Time      `json:"createdAt"`
	Choices     []AnswerChoice `json:"choices"`
}

type AnswerChoice struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type AnswerResult struct {
	IsCorrect       bool      `json:"isCorrect"`
	UserAnswer      string    `json:"userAnswer"`
	CorrectAnswer   string    `json:"correctAnswer"`
	FeedbackTitle   string    `json:"feedbackTitle"`
	FeedbackMessage string    `json:"feedbackMessage"`
	RedFlags        []RedFlag `json:"redFlags"`
	Explanation     string    `json:"explanation"`
}
