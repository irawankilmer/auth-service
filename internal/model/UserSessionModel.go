package model

import "time"

type UserSession struct {
	ID               string
	UserID           string
	RefreshTokenHash string
	DeviceID         string
	IPAddress        string
	UserAgent        string
	Revoked          bool
	ExpiresAt        time.Time
}
