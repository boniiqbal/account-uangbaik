package http

import (
	_config "uangbaik-account-microservice/config"
	"uangbaik-account-microservice/database/models"

	"github.com/gin-gonic/gin"
)

//ChangeBank payload struct
type EditUser struct {
	Data struct {
		models.User
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

// EdirUserValidator .
func EditUserValidator() EditUser {
	editUserValidator := EditUser{}
	return editUserValidator
}
