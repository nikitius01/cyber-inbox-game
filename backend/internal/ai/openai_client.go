package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"cybersecurity-game/backend/internal/tasks"
)

type Service struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
	mu      sync.RWMutex
	tasks   map[string]tasks.Task
}

func NewService(apiKey, baseURL, model, timeoutSeconds string) *Service {
	timeout := 90 * time.Second
	if parsed, err := strconv.Atoi(timeoutSeconds); err == nil && parsed >= 10 && parsed <= 300 {
		timeout = time.Duration(parsed) * time.Second
	}
	return &Service{
		apiKey:  apiKey,
		baseURL: normalizeChatCompletionsURL(baseURL),
		model:   model,
		client:  &http.Client{Timeout: timeout},
		tasks:   map[string]tasks.Task{},
	}
}

func normalizeChatCompletionsURL(baseURL string) string {
	baseURL = strings.TrimRight(baseURL, "/")
	if strings.HasSuffix(baseURL, "/chat/completions") {
		return baseURL
	}
	if strings.HasSuffix(baseURL, "/v1") {
		return baseURL + "/chat/completions"
	}
	return baseURL
}

func (s *Service) GenerateTask() (tasks.Task, error) {
	if s.apiKey == "" {
		return tasks.Task{}, errors.New("AI_API_KEY or OPENAI_API_KEY is not configured")
	}
	reqBody := map[string]any{
		"model": s.model,
		"messages": []map[string]string{
			{"role": "system", "content": Prompt},
			{"role": "user", "content": "Сгенерируй одну задачу категории AI для игры Инспектор входящих."},
		},
		"response_format": map[string]string{"type": "json_object"},
	}
	raw, _ := json.Marshal(reqBody)
	req, err := http.NewRequest(http.MethodPost, s.baseURL, bytes.NewReader(raw))
	if err != nil {
		return tasks.Task{}, err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	if err != nil {
		return tasks.Task{}, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return tasks.Task{}, errors.New("AI provider request failed")
	}
	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return tasks.Task{}, err
	}
	if len(out.Choices) == 0 {
		return tasks.Task{}, errors.New("OpenAI returned no choices")
	}
	var task tasks.Task
	if err := json.Unmarshal([]byte(out.Choices[0].Message.Content), &task); err != nil {
		return tasks.Task{}, err
	}
	task.Category = tasks.AI
	if err := validateTask(task); err != nil {
		return tasks.Task{}, err
	}
	s.mu.Lock()
	s.tasks[task.ID] = task
	s.mu.Unlock()
	return task, nil
}

func (s *Service) CheckAnswer(r *http.Request) (tasks.AnswerResult, error) {
	var req struct {
		TaskID string `json:"taskId"`
		Answer string `json:"answer"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return tasks.AnswerResult{}, err
	}
	s.mu.RLock()
	task, ok := s.tasks[req.TaskID]
	s.mu.RUnlock()
	if !ok {
		return tasks.AnswerResult{}, errors.New("AI task not found")
	}
	correct := "legitimate"
	if task.IsPhishing {
		correct = "phishing"
	}
	return tasks.BuildAnswerResult(task, req.Answer, correct), nil
}

func validateTask(task tasks.Task) error {
	if task.ID == "" || task.Subject == "" || len(task.Body) > 5000 || len(task.Raw.Source) > 12000 {
		return errors.New("invalid AI task")
	}
	return nil
}
