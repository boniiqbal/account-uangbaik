package http

import (
	_config "uangbaik-account-microservice/config"

	"github.com/gin-gonic/gin"
)

//Pin payload struct
type Pin struct {
	Data struct {
		PinUser string `form:"pin_user" json:"pin_user" binding:"required,min=6"`
	} `json:"data"`
}

//BindCreatePin for bind request
func (pin *Pin) BindCreatePin(c *gin.Context) error {
	err := _config.BindRequest(c, pin)
	if err != nil {
		return err
	}
	return nil
}

// PinValidator .
func PinValidator() Pin {
	pinValidator := Pin{}
	return pinValidator
}
