package main

import (
	"database/sql"
	"log"
	"net/http"

	"cybersecurity-game/backend/internal/ai"
	"cybersecurity-game/backend/internal/auth"
	"cybersecurity-game/backend/internal/config"
	"cybersecurity-game/backend/internal/middleware"
	"cybersecurity-game/backend/internal/profile"
	"cybersecurity-game/backend/internal/tasks"
	"cybersecurity-game/backend/internal/users"
)

func main() {
	cfg := config.Load()

	var userRepo users.Repository
	if cfg.DatabaseURL != "" {
		db, err := sql.Open("postgres", cfg.DatabaseURL)
		if err != nil {
			log.Fatal(err)
		}
		if err := db.Ping(); err != nil {
			log.Fatal(err)
		}
		postgresRepo := users.NewPostgresRepository(db)
		if err := postgresRepo.Migrate(); err != nil {
			log.Fatal(err)
		}
		userRepo = postgresRepo
		log.Println("users repository: PostgreSQL")
	} else {
		userRepo = users.NewMemoryRepository()
		log.Println("users repository: in-memory")
	}
	taskRepo := tasks.NewMemoryRepository()
	taskRepo.Seed(tasks.BuildSeedTasks())

	authService := auth.NewService(userRepo, cfg.JWTSecret)
	taskService := tasks.NewService(taskRepo)
	profileService := profile.NewService(userRepo)
	aiService := ai.NewService(cfg.AIAPIKey, cfg.AIBaseURL, cfg.AIModel, cfg.AITimeoutSeconds)

	mux := http.NewServeMux()
	auth.RegisterRoutes(mux, authService)
	users.RegisterRoutes(mux, userRepo, cfg.EnableDevUserList)
	tasks.RegisterRoutes(mux, taskService, aiService, userRepo)
	profile.RegisterRoutes(mux, profileService)

	handler := middleware.Chain(
		mux,
		middleware.SecurityHeaders,
		middleware.CORS(cfg.FrontendOrigin),
		middleware.RateLimit(),
		middleware.AuthContext(authService),
	)

	log.Printf("server listening on %s", cfg.Address)
	if err := http.ListenAndServe(cfg.Address, handler); err != nil {
		log.Fatal(err)
	}
}
