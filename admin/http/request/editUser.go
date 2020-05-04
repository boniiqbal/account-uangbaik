package http

import (
	_config "uangbaik-account-microservice/config"

	"github.com/gin-gonic/gin"
)

//EditUser payload struct
type EditUser struct {
	Data struct {
		FullName        string `form:"full_name" json:"full_name"`
		Phone           string `form:"phone" json:"phone"`
		Email           string `form:"email" json:"email"`
		Address         string `form:"address" json:"address"`
		Status          int    `form:"status" json:"status"`
		IsPhoneVerified int    `form:"is_phone_verified" json:"is_phone_verified"`
		IsEmailVerified int    `form:"is_email_verified" json:"is_email_verified"`
	} `json:"data"`
}

//BindEditUser for bind request
func (user *EditUser) BindEditUser(c *gin.Context) error {
	err := _config.BindRequest(c, user)
	if err != nil {
		return err
	}
	return nil
}

// EditUserValidator .
func EditUserValidator() EditUser {
	editUserValidator := EditUser{}
	return editUserValidator
}
