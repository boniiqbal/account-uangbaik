package models

import (
	"time"

	"github.com/google/uuid"
)

//User model
type User struct {
	ID              uint          `gorm:"primary_key" json:"id"`
	UUID            uuid.UUID     `json:"uuid"`
	FullName        string        `gorm:"size:255" json:"full_name"`
	Phone           string        `gorm:"size:255;unique_index;not null" json:"phone"`
	Email           string        `gorm:"size:255" json:"email"`
	Address         string        `gorm:"size:255" json:"address"`
	Pin             string        `gorm:"size:255" json:"-"`
	Password        string        `gorm:"size:255" json:"password"`
	Hash            string        `gorm:"size:255" json:"hash"`
	LastLogin       *time.Time    `json:"last_login"`
	Status          int           `gorm:"default:0" json:"status"`
	IsPhoneVerified int           `gorm:"default:0" json:"is_phone_verified"`
	IsEmailVerified int           `gorm:"default:0" json:"is_email_verified"`
	BankAccount     []BankAccount `gorm:"polymorphic:Actor;" json:"bank_account"`
	AccessToken     []AccessToken `json:"access_token"`
	UploadAsset     []UploadAsset `gorm:"polymorphic:Upload;" json:"upload_assets"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	DeletedAt       *time.Time    `json:"deleted_at"`
}

//BeforeCreate for create uuid
func (u *User) BeforeCreate() (err error) {
	u.UUID = uuid.New()
	return
}
