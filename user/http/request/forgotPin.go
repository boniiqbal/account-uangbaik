package http

import (
	_config "uangbaik-account-microservice/config"

	"github.com/gin-gonic/gin"
)

//ForgotPin payload struct
type ForgotPin struct {
	Data struct {
		NewPin string `form:"new_pin" json:"new_pin" binding:"required,min=6"`
	} `json:"data"`
}

//BindForgotPin for bind request
func (forgot *ForgotPin) BindForgotPin(c *gin.Context) error {
	err := _config.BindRequest(c, forgot)
	if err != nil {
		return err
	}
	return nil
}

// ForgotPinValidator .
func ForgotPinValidator() ForgotPin {
	forgotPinValidator := ForgotPin{}
	return forgotPinValidator
}
