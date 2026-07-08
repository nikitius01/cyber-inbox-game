package config

import "os"

type Config struct {
	Address           string
	FrontendOrigin    string
	JWTSecret         string
	AIAPIKey          string
	AIBaseURL         string
	AIModel           string
	AITimeoutSeconds  string
	DatabaseURL       string
	EnableDevUserList bool
}

func Load() Config {
	return Config{
		Address:           getEnv("SERVER_ADDR", ":8080"),
		FrontendOrigin:    getEnv("FRONTEND_ORIGIN", "http://localhost:5173"),
		JWTSecret:         getEnv("JWT_SECRET", "dev-change-me-long-random-secret"),
		AIAPIKey:          getEnv("AI_API_KEY", os.Getenv("OPENAI_API_KEY")),
		AIBaseURL:         getEnv("AI_BASE_URL", "https://api.openai.com/v1/chat/completions"),
		AIModel:           getEnv("AI_MODEL", "gpt-4.1-mini"),
		AITimeoutSeconds:  getEnv("AI_TIMEOUT_SECONDS", "90"),
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		EnableDevUserList: os.Getenv("ENABLE_DEV_USER_LIST") == "true",
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
