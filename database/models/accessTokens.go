package models

import "time"

//AccessToken for model access token
type AccessToken struct {
	ID           uint       `gorm:"primary_key" json:"id"`
	UserID       uint       `gorm:"not null" json:"user_id"`
	AccessToken  string     `gorm:"size:255;not null" json:"access_token"`
	TokenType    string     `json:"token_type"`
	Scope        string     `json:"scope"`
	RefreshToken string     `gorm:"size:255" json:"refresh_token"`
	UserAgent    string     `json:"user_agent"`
	ClientIP     string     `gorm:"size:255" json:"client_ip"`
	ExpiresAt    int64      `json:"expires_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}
