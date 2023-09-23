package model

import "time"

type Session struct {
	RefreshToken string    `json:"refreshToken" bson:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt" bson:"expiresAt"`
	IssuedAt     time.Time `json:"issuedAt" bson:"issuedAt"`
}
