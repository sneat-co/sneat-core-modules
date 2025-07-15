package dbo4invitus

import "time"

// InviteClaim record
type InviteClaim struct {
	Time   time.Time `json:"time" firestore:"time"`
	UserID string    `json:"userId" firestore:"userId"`
}

// InviteCode record
type InviteCode struct {
	Claims []InviteClaim `json:"claims,omitempty" firestore:"claims,omitempty"`
}
