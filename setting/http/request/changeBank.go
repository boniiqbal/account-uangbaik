package http

import (
	_config "uangbaik-account-microservice/config"
	"uangbaik-account-microservice/database/models"

	"github.com/gin-gonic/gin"
)

//ChangeBank payload struct
type ChangeBank struct {
	Data struct {
		models.BankAccount
	} `json:"data"`
}

//BindCreateBank for bind request
func (bank *ChangeBank) BindChangeBank(c *gin.Context) error {
	err := _config.BindRequest(c, bank)
	if err != nil {
		return err
	}
	return nil
}

// ChangeBankValidator .
func ChangeBankValidator() ChangeBank {
	changeBankValidator := ChangeBank{}
	return changeBankValidator
}
