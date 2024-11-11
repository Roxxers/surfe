package services

import (
	"math"
	"sort"

	"github.com/roxxers/surfe-techtest/internal/core/domain"
	"github.com/roxxers/surfe-techtest/internal/ports"
)

type Service struct {
	db ports.Database
}

func NewService(db ports.Database) *Service {
	return &Service{db: db}
}

func (s *Service) FetchUser(userId int64) *domain.User {
	user := s.db.GetUser(userId)
	return &domain.User{
		Id:        user.Id,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
}

func (s *Service) GetUserActionCount(userId int64) int32 {
	actions := s.db.GetActionsForUser(userId)
	return int32(len(actions))
}

// CalculateNextActionProbablity checks all occurances of next action by all users and returns a map of actions and their probablity of being the next action taken by a user.
// This is done via looking at what the total amount of that action is in proportion to the total number of actions taken after a given action.
func (s *Service) CalculateNextActionProbablity(actionType string) map[string]float64 {
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
				if _, ok := nextActions[action.Type]; !ok {
					nextActions[action.Type] = 0
				}
				nextActions[actions[i+1].Type]++
				totalActions++
			}
		}
	}

	for action, count := range nextActions {
		probablity := float64(count) / float64(totalActions)
		// This is to make the end result look like the API Example, but it could lead to the total not equalling 1 which would make the probablities slightly off.
		actionProbablities[action] = math.Round(probablity*100) / 100
	}
	return actionProbablities
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
			if action.Type == "REFER_USER" {
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
