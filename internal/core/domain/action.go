package domain

import "time"

type Action struct {
	Id           uint64
	Type         string
	UserId       uint64
	TargetUserId uint64
	CreatedAt    time.Time
}
