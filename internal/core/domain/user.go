package domain

import "time"

type User struct {
	Id        uint64
	Name      string
	CreatedAt time.Time
}
