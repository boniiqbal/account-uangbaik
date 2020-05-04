package http

import (
	_config "uangbaik-account-microservice/config"

	"github.com/gin-gonic/gin"
)

//User payload struct
type User struct {
	Data struct {
		Phone string `form:"phone" json:"phone" binding:"required"`
	} `json:"data"`
}

//Bind for bind request
func (user *User) Bind(c *gin.Context) error {
	err := _config.BindRequest(c, user)
	if err != nil {
		return err
	}
	return nil
}

// UserValidator .
func UserValidator() User {
	userValidator := User{}
	return userValidator
}
