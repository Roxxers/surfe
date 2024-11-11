package secondary

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/roxxers/surfe-techtest/internal/core/domain"
)

var (
	ErrNoSuchUserID     = errors.New("user with that id does not exist")
	ErrNoSuchActionID   = errors.New("action with that id does not exist")
	ErrNoSuchActionType = errors.New("action with that type does not exist")
)

type UserTable = map[int64]domain.User
type ActionTable = map[int64]domain.Action
type ActionUserIDIndex = map[int64][]domain.Action

// ActionsUserIDIndex Is a simulation of indexing actions on the UserID of the action. We will use this for later functions.

type MemoryDatabase struct {
	Users              UserTable
	Actions            ActionTable
	ActionsUserIDIndex ActionUserIDIndex
}

func (db *MemoryDatabase) GetUser(userID int64) (domain.User, error) {
	if _, ok := db.Users[userID]; !ok {
		return domain.User{}, ErrNoSuchUserID
	}
	return db.Users[userID], nil
}

func (db *MemoryDatabase) GetActionsForUser(userID int64) ([]domain.Action, error) {
	if _, ok := db.ActionsUserIDIndex[userID]; !ok {
		return nil, ErrNoSuchUserID
	}
	return db.ActionsUserIDIndex[userID], nil
}

func (db *MemoryDatabase) GetAction(actionId int64) (domain.Action, error) {
	if _, ok := db.Actions[actionId]; !ok {
		return domain.Action{}, ErrNoSuchActionID
	}
	return db.Actions[actionId], nil
}

func (db *MemoryDatabase) GetAllUsers() map[int64]domain.User {
	return db.Users
}

func (db *MemoryDatabase) GetActionsPerUser() map[int64][]domain.Action {
	return db.ActionsUserIDIndex
}

func NewMemoryDatabase() *MemoryDatabase {
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
	file, err := os.Open("./internal/adapters/secondary/users.json") // Hardcoded for ease
	defer file.Close()
	if err != nil {
		return nil, err
	}
	usersByte, _ := io.ReadAll(file)
	var users []domain.User
	json.Unmarshal(usersByte, &users)

	// O(n) on number of users but we are doing it at startup
	// Would not need to do this for an actual inProd program
	userTable := make(UserTable)
	for _, user := range users {
		userTable[user.Id] = user
	}
	users = nil
	return userTable, nil
}

// Mostly copied pasted function, could make it generic but no need to do so for this test.
func loadActions() (ActionTable, ActionUserIDIndex, error) {
	file, err := os.Open("./internal/adapters/secondary/actions.json")
	defer file.Close()
	if err != nil {
		return nil, nil, err
	}
	actionByte, _ := io.ReadAll(file)
	var actions []domain.Action
	json.Unmarshal(actionByte, &actions)

	// O(n) on number of users but we are doing it at startup
	// Would not need to do this for an actual inProd program
	actionTable := make(ActionTable)
	for _, action := range actions {
		actionTable[action.Id] = action
	}

	userIDIndex := createActionsUserIdIndex(actions)
	actions = nil
	return actionTable, userIDIndex, nil
}

func createActionsUserIdIndex(actions []domain.Action) ActionUserIDIndex {
	actionsUserIDIndex := make(ActionUserIDIndex)
	for _, action := range actions {
		actionsUserIDIndex[action.UserId] = append(actionsUserIDIndex[action.UserId], action)
	}
	return actionsUserIDIndex
}
