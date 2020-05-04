package http

import (
	_config "uangbaik-account-microservice/config"

	"github.com/gin-gonic/gin"
)

//BankAccount payload struct
type BankAccount struct {
	Data struct {
		BankID        uint   `form:"bank_id" json:"bank_id"  binding:"required"`
		AccountNumber string `form:"account_number" json:"account_number" binding:"required"`
	} `json:"data"`
}

//BindCreateBank for bind request
func (bank *BankAccount) BindCreateBank(c *gin.Context) error {
	err := _config.BindRequest(c, bank)
	if err != nil {
		return err
	}
	return nil
}

// BankValidator .
func BankValidator() BankAccount {
	bankValidator := BankAccount{}
	return bankValidator
}
