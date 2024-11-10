package ports

import "github.com/roxxers/surfe-techtest/internal/core/domain"

type Database interface {
	GetAction(actionId uint64) domain.Action
	GetActionsForUser(userId uint64) []domain.Action
	GetUser(userID uint64) domain.User
}
