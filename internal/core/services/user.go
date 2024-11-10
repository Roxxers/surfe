package services

import (
	"github.com/roxxers/surfe-techtest/internal/core/domain"
	"github.com/roxxers/surfe-techtest/internal/ports"
)

type Service struct {
	db ports.Database
}

func NewService(db ports.Database) *Service {
	return &Service{db: db}
}

func (s *Service) FetchUser(userId uint64) *domain.User {
	user := s.db.GetUser(userId)
	return &domain.User{
		Id:        user.Id,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
}

func (s *Service) GetUserActionCount(userId uint64) int32 {
	actions := s.db.GetActionsForUser(userId)
	return int32(len(actions))
}
