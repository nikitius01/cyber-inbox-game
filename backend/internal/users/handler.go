package users

import (
	"encoding/json"
	"net/http"

	"cybersecurity-game/backend/internal/middleware"
)

func RegisterRoutes(mux *http.ServeMux, repo Repository, enableDevUserList bool) {
	mux.HandleFunc("GET /api/auth/me", func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := repo.FindByID(userID)
		if err != nil {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}
		writeJSON(w, http.StatusOK, user)
	})

	mux.HandleFunc("GET /api/dev/users", func(w http.ResponseWriter, r *http.Request) {
		if !enableDevUserList {
			http.Error(w, "dev user list is disabled", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, repo.List())
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
