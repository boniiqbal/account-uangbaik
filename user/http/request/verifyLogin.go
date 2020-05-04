package http

import (
	_config "uangbaik-account-microservice/config"

	"github.com/gin-gonic/gin"
)

//Verify payload struct
type Verify struct {
	Data struct {
		Phone string `form:"phone" json:"phone" binding:"required"`
		Pin   string `form:"pin" json:"pin" binding:"required,min=6"`
	} `json:"data"`
}

//BindVerify for bind request
func (verify *Verify) BindVerify(c *gin.Context) error {
	err := _config.BindRequest(c, verify)
	if err != nil {
		return err
	}
	return nil
}

// VerifyValidator .
func VerifyValidator() Verify {
	verifyValidator := Verify{}
	return verifyValidator
}
