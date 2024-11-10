package ports

import "github.com/roxxers/surfe-techtest/internal/core/domain"

type Database interface {
	GetAction(actionId int64) domain.Action
	GetActionsForUser(userId int64) []domain.Action
	GetUser(userID int64) domain.User
}
