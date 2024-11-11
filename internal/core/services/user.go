package services

import (
	"errors"
	"math"
	"sort"

	"github.com/roxxers/surfe-techtest/internal/core/domain"
	"github.com/roxxers/surfe-techtest/internal/ports"
)

// In a larger program this would hopefully be somewhere easier to manage like a file of constants shared for the program.
// For this example, it is only used here so we will define it here.
const (
	REFER_USER_ACTION = "REFER_USER"
)

type Service struct {
	db ports.Database
}

func NewService(db ports.Database) *Service {
	return &Service{db: db}
}

// Fetch user returns a domain.User from storage
func (s *Service) FetchUser(userId int64) (*domain.User, error) {
	user, err := s.db.GetUser(userId)
	if err != nil {
		return &domain.User{}, err
	}
	return &domain.User{
		Id:        user.Id,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}

// GetUserActionCount returns the total number of actions for the specified user
func (s *Service) GetUserActionCount(userId int64) (int32, error) {
	actions, err := s.db.GetActionsForUser(userId)
	if err != nil {
		return 0, err
	}
	return int32(len(actions)), nil
}

// CalculateNextActionProbablity checks all occurances of next action by all users and returns a map of actions and their probablity of being the next action taken by a user.
// This is done via looking at what the total amount of that action is in proportion to the total number of actions taken after a given action.
func (s *Service) CalculateNextActionProbablity(actionType string) (map[string]float64, error) {
	// Since the way the task is worded does not require me to make this function as performant as possible, I will be sorting the lists in the function
	// In a more optimal world, we would have these sorted at a database level, with either making those calls in the database layer (to save on memory usage here)
	// or with k-orderable ID's which would allow us to make assuptions when our ID's come through to this side if the database layer is sharded.
	actionsPerUser := s.db.GetActionsPerUser()
	nextActions := make(map[string]int)
	totalActions := 0
	actionProbablities := make(map[string]float64)

	for _, actions := range actionsPerUser {
		sort.Slice(actions, func(i, j int) bool {
			return actions[i].CreatedAt.Before(actions[j].CreatedAt)
		})
		for i, action := range actions {
			if action.Type == actionType {
				if i == len(actions)-1 {
					// This is the last action of the user, cannot check the next one
					continue
				}
				nextAction := actions[i+1]
				if _, ok := nextActions[nextAction.Type]; !ok {
					nextActions[nextAction.Type] = 0
				}
				nextActions[nextAction.Type]++
				totalActions++
			}
		}
	}

	if len(nextActions) == 0 {
		// Action wasn't found or always at the end of an action chain - very unlikely!!
		return nil, errors.New("could not find specified type - or very unlikely the final event for all users")
	}

	for action, count := range nextActions {
		probablity := float64(count) / float64(totalActions)
		// This is to make the end result look like the API Example, but it could lead to the total not equalling 1 which would make the probablities slightly off.
		actionProbablities[action] = math.Round(probablity*100) / 100
	}

	return actionProbablities, nil
}

// Designing this function keeping in mind complexity, it is only limited to the size of the users and all of their actions,
// We avoid pitfalls of looking for actions per user via the indexing, the hashmap makes any call to it O(1).
// Functions time complexity is O(nm) (size of users and size of all actions). We only loop over once and only values we need to consider.
// This is helped with our pre indexed actions list which is a hashmap of actions grouped per user. This simulates the work a
// database would do in opitmizing datasets before we have them in this service.
// In practice, O(nm) is the worst case (where every user is referred to from one base user and every action is a referal).
// We optimize some running through every single user by keeping track of all previously seen users and calculating their value during the recursion loop.
// Space complexity will be O(n) on the size of users as we keep a hashmap that is as long as the userset. The recursion stack can also get to O(n) if every user is traced back to a single user.

// CalculateAllUserReferalIndexes calculates the amount of referred users that can be tied back to an originating user.
// Calculating a total number of users directly or indirectly (via the users they directly referred) referred to the platform by a particular user.
// Returns map[int64]int mapping user id's to their "Referal Index"
func (s *Service) CalculateAllUserReferalIndexes() map[int64]int {
	actionsPerUser := s.db.GetActionsPerUser()
	ReferalIndex := make(map[int64]int)

	// treeSearching is a recursive function,
	// Returns int representing the referal index of the node (user)
	// that its calculating (all users referred to the service below that current user)
	var treeSearching func(userId int64) int
	treeSearching = func(userId int64) int {
		if _, ok := ReferalIndex[userId]; ok {
			return ReferalIndex[userId] // We have already calculated this
		}
		childCount := 0

		for _, action := range actionsPerUser[userId] {
			if action.Type == REFER_USER_ACTION {
				if userId != action.TargetUserId {
					// Dataset includes actions with users that refer themselves.
					// 802 being a prime example
					// This seems wrong so we will ignore this
					childCount++
					returncc := treeSearching(action.TargetUserId)
					childCount += returncc
				}
			}
		}
		ReferalIndex[userId] = childCount
		return childCount
	}

	for userId, _ := range actionsPerUser {
		treeSearching(userId)
	}
	return ReferalIndex
}
