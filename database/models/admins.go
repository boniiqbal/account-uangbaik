package models

import "time"

//Admin for model access token
type Admin struct {
	ID              uint               `gorm:"primary_key" json:"id"`
	FullName        string             `gorm:"size:255" json:"full_name"`
	Phone           string             `gorm:"size:255;unique_index;not null" json:"phone"`
	Email           string             `gorm:"size:255" json:"email"`
	Address         string             `gorm:"size:255" json:"address"`
	Username        string             `gorm:"size:255" json:"username"`
	Password        string             `gorm:"size:255" json:"-"`
	Hash            string             `gorm:"size:255" json:"hash"`
	LastLogin       *time.Time         `json:"last_login"`
	Status          int                `gorm:"default:0" json:"status"`
	AdminAccesToken []AdminAccessToken `json:"admin_access_token"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	DeletedAt       *time.Time         `json:"deleted_at"`
}
