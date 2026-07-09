package tasks

import (
	"errors"
	"math/rand"
)

var ErrBadAnswer = errors.New("answer must be phishing or legitimate")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RandomTask(category Difficulty) (PublicTask, error) {
	task, err := s.repo.RandomByCategory(category)
	if err != nil {
		return PublicTask{}, err
	}
	return ToPublicTask(task), nil
}

func (s *Service) CheckAnswer(id, answer string) (Task, AnswerResult, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return Task{}, AnswerResult{}, err
	}
	if answer != "phishing" && answer != "legitimate" {
		return Task{}, AnswerResult{}, ErrBadAnswer
	}
	correct := "legitimate"
	if task.IsPhishing {
		correct = "phishing"
	}
	return task, BuildAnswerResult(task, answer, correct), nil
}

func BuildAnswerResult(task Task, userAnswer, correct string) AnswerResult {
	isCorrect := userAnswer == correct
	title := "Верный разбор"
	message := "Ответ принят. Ты правильно определил тип письма."
	if !isCorrect {
		title = "Почему ответ неверный"
		message = "Нужно было выбрать: " + answerLabel(correct) + ". " + task.Explanation
	}
	return AnswerResult{
		IsCorrect:       isCorrect,
		UserAnswer:      userAnswer,
		CorrectAnswer:   correct,
		FeedbackTitle:   title,
		FeedbackMessage: message,
		RedFlags:        task.RedFlags,
		Explanation:     task.Explanation,
	}
}

func answerLabel(answer string) string {
	if answer == "phishing" {
		return "Фишинг"
	}
	return "Легитимное письмо"
}

func ToPublicTask(task Task) PublicTask {
	choices := []AnswerChoice{
		{Label: "Фишинг", Value: "phishing"},
		{Label: "Легитимное", Value: "legitimate"},
	}
	rand.Shuffle(len(choices), func(i, j int) { choices[i], choices[j] = choices[j], choices[i] })
	return PublicTask{
		ID:          task.ID,
		Category:    task.Category,
		SenderName:  task.SenderName,
		SenderEmail: task.SenderEmail,
		Subject:     task.Subject,
		Body:        task.Body,
		Links:       task.Links,
		Attachments: task.Attachments,
		Raw:         task.Raw,
		CreatedAt:   task.CreatedAt,
		Choices:     choices,
	}
}
