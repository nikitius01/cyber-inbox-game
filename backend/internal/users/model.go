package users

import "time"

type Stats struct {
	Streak           int            `json:"streak"`
	TotalSolved      int            `json:"totalSolved"`
	SolvedByCategory map[string]int `json:"solvedByCategory"`
	LastActivityDate string         `json:"lastActivityDate"`
}

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Stats        Stats     `json:"stats"`
	CreatedAt    time.Time `json:"createdAt"`
}
