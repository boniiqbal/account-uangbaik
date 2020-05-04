package models

import (
	"errors"
	"time"
)

//BankAccount for model bank account
type BankAccount struct {
	ID            uint       `gorm:"primary_key" json:"id"`
	ActorID       uint       `json:"actor_id"`
	ActorType     string     `gorm:"size:45" json:"actor_type"`
	BankName      string     `gorm:"size:255" json:"bank_name"`
	BankBranch    string     `gorm:"size:255" json:"bank_branch"`
	AccountNumber string     `gorm:"size:255" json:"account_number"`
	Status        int        `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

//BeforeSave for 
func (u *BankAccount) BeforeSave() (err error) {
	if u.BankName == "" || u.BankBranch == "" || u.AccountNumber == "" {
		err = errors.New("Data is required")
	}
	return err
}
