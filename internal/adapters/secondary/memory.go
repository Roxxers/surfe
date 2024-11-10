package secondary

import (
	"encoding/json"
	"io"
	"os"

	"github.com/roxxers/surfe-techtest/internal/core/domain"
)

type UserTable = map[uint64]domain.User
type ActionTable = map[uint64]domain.Action
type ActionUserIDIndex = map[uint64][]domain.Action

// ActionsUserIDIndex Is a simulation of indexing actions on the UserID of the action. We will use this for later functions.

type MemoryDatabase struct {
	Users              UserTable
	Actions            ActionTable
	ActionsUserIDIndex ActionUserIDIndex
}

func (db *MemoryDatabase) GetUser(userID uint64) domain.User {
	return db.Users[userID]
}

func (db *MemoryDatabase) GetActionsForUser(userID uint64) []domain.Action {
	return db.ActionsUserIDIndex[userID]
}

func (db *MemoryDatabase) GetAction(actionId uint64) domain.Action {
	return db.Actions[actionId]
}

func LoadMemoryDatabase() *MemoryDatabase {
	userTable, err := loadUsers()
	if err != nil {
		panic(err)
	}

	actionTable, actionUserIDIndex, err := loadActions()
	if err != nil {
		panic(err)
	}

	return &MemoryDatabase{
		Users:              userTable,
		Actions:            actionTable,
		ActionsUserIDIndex: actionUserIDIndex,
	}
}

func loadUsers() (UserTable, error) {
	file, err := os.Open("./users.json")
	defer file.Close()
	if err != nil {
		return nil, err
	}
	usersByte, _ := io.ReadAll(file)
	var users []domain.User
	json.Unmarshal(usersByte, &users)

	// O(n) on number of users but we are doing it at startup
	// Would not need to do this for an actual inProd program
	var userTable UserTable
	for _, user := range users {
		userTable[user.Id] = user
	}
	users = nil
	return userTable, nil
}

// Mostly copied pasted function, could make it generic but no need to do so for this test.
func loadActions() (ActionTable, ActionUserIDIndex, error) {
	file, err := os.Open("./action.json")
	defer file.Close()
	if err != nil {
		return nil, nil, err
	}
	actionByte, _ := io.ReadAll(file)
	var actions []domain.Action
	json.Unmarshal(actionByte, &actions)

	// O(n) on number of users but we are doing it at startup
	// Would not need to do this for an actual inProd program
	var actionTable ActionTable
	for _, action := range actions {
		actionTable[action.Id] = action
	}

	userIDIndex := createActionsUserIdIndex(actions)
	actions = nil
	return actionTable, userIDIndex, nil
}

func createActionsUserIdIndex(actions []domain.Action) ActionUserIDIndex {
	var actionsUserIDIndex ActionUserIDIndex
	for _, action := range actions {
		actionsUserIDIndex[action.UserId] = append(actionsUserIDIndex[action.UserId], action)
	}
	return actionsUserIDIndex
}
