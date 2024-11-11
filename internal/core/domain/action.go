package domain

import "time"

// Defining the json mapping here does couple the adapter with the domain and I'm doing it here to save time

type Action struct {
	Id           int64     `json:"id"`
	Type         string    `json:"type"` // Possibly Enum here for better type handling when all possible actions are defined
	UserId       int64     `json:"userId"`
	TargetUserId int64     `json:"targetUser,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}
