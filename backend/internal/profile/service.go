package profile

import "cybersecurity-game/backend/internal/users"

type Service struct {
	users users.Repository
}

func NewService(repo users.Repository) *Service {
	return &Service{users: repo}
}

func (s *Service) Profile(userID string) (users.User, error) {
	return s.users.FindByID(userID)
}
