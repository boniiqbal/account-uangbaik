package models

import (
	"errors"
	"time"
)

//Bank for model bank account
type Bank struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

//BeforeCreate for
func (u *Bank) BeforeCreate() (err error) {
	if u.Code == "" {
		err = errors.New("Code is required")
	}
	if u.Name == "" {
		err = errors.New("Name is required")
	}

	return err
}
