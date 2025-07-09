package model

import "time"

type EmailVerificationModel struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	IsUsed    bool
	CreatedAt time.Time
}
