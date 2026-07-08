package tasks

import (
	"encoding/json"
	"net/http"
	"time"

	"cybersecurity-game/backend/internal/middleware"
	"cybersecurity-game/backend/internal/users"
)

type answerRequest struct {
	Answer string `json:"answer"`
}

type AIService interface {
	GenerateTask() (Task, error)
	CheckAnswer(*http.Request) (AnswerResult, error)
}

func RegisterRoutes(mux *http.ServeMux, service *Service, aiService AIService, userRepo users.Repository) {
	mux.HandleFunc("GET /api/tasks/random", func(w http.ResponseWriter, r *http.Request) {
		category := Difficulty(r.URL.Query().Get("category"))
		if category == "" {
			category = Easy
		}
		if category == AI {
			task, err := aiService.GenerateTask()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
				return
			}
			writeJSON(w, http.StatusOK, ToPublicTask(task))
			return
		}
		task, err := service.RandomTask(category)
		if err != nil {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, task)
	})

	mux.HandleFunc("POST /api/tasks/{id}/answer", func(w http.ResponseWriter, r *http.Request) {
		var req answerRequest
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		task, result, err := service.CheckAnswer(r.PathValue("id"), req.Answer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		updateUserStats(r, userRepo, string(task.Category))
		writeJSON(w, http.StatusOK, result)
	})

	mux.HandleFunc("POST /api/ai/tasks/answer", func(w http.ResponseWriter, r *http.Request) {
		result, err := aiService.CheckAnswer(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		updateUserStats(r, userRepo, string(AI))
		writeJSON(w, http.StatusOK, result)
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func updateUserStats(r *http.Request, repo users.Repository, category string) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		return
	}
	user, err := repo.FindByID(userID)
	if err != nil {
		return
	}
	if user.Stats.SolvedByCategory == nil {
		user.Stats.SolvedByCategory = map[string]int{}
	}
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	if user.Stats.LastActivityDate == "" {
		user.Stats.Streak = 1
	} else if user.Stats.LastActivityDate == yesterday {
		user.Stats.Streak++
	} else if user.Stats.LastActivityDate != today {
		user.Stats.Streak = 1
	}
	user.Stats.TotalSolved++
	user.Stats.SolvedByCategory[category]++
	user.Stats.LastActivityDate = today
	_ = repo.Update(user)
}
