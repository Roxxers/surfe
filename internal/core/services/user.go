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

// Designing this function for complexity, our time complexity is only limited to the size of the users and all of their actions,
// We avoid pitfalls of looking for actions per user via the indexing, the hashmap makes any call to it O(1).
// Function should run O(nm) (size of users and size of all actions). We only loop over once and only values we need to consider.
// This is helped with our pre indexed actions list which created a hashmap of actions based off their user.
// In the real world, this runs faster than O(nm) as it is only O(nm) in the worst case, which would be 0 referred users.
// There is saved time by building the referal index of users invited as re go through.
func (s *Service) CalculateAllUserReferalIndexes() map[int64]int {
	actionsPerUser := s.db.GetActionsPerUser()
	ReferalIndex := make(map[int64]int)

	// treeSearching is a recursive function, returns int representing the referal index of the node that its being run on
	var treeSearching func(userId int64) int
	treeSearching = func(userId int64) int {
		if _, ok := ReferalIndex[userId]; ok {
			return ReferalIndex[userId] // We have already calculated this
		}
		childCount := 0

		for _, action := range actionsPerUser[userId] {
			if action.Type == "REFER_USER" {
				if userId != action.TargetUserId {
					childCount++
					// Dataset includes actions with users that refer themselves. This seems wrong so we will ignore this and not even count it
					// 802 being a prime example
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
