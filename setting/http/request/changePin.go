package http

import (
	_config "uangbaik-account-microservice/config"

	"github.com/gin-gonic/gin"
)

//ChangePin payload struct
type ChangePin struct {
	Data struct {
		OldPin      string `form:"old_pin" json:"old_pin"  binding:"required"`
		NewPin    string `form:"new_pin" json:"new_pin" binding:"required"`
	} `json:"data"`
}

//BindChangePin for bind request
func (pin *ChangePin) BindChangePin(c *gin.Context) error {
	err := _config.BindRequest(c, pin)
	if err != nil {
		return err
	}
	return nil
}

// ChangePinValidator .
func ChangePinValidator() ChangePin {
	changePinValidator := ChangePin{}
	return changePinValidator
}
