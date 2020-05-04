package http

import (
	_config "uangbaik-account-microservice/config"

	"github.com/gin-gonic/gin"
)

//LoginAdmin payload struct
type LoginAdmin struct {
	Data struct {
		Username string `form:"username" json:"username" binding:"required"`
		Password   string `form:"password" json:"password" binding:"required"`
	} `json:"data"`
}

//BindLoginAdmin for bind request
func (login *LoginAdmin) BindLoginAdmin(c *gin.Context) error {
	err := _config.BindRequest(c, login)
	if err != nil {
		return err
	}
	return nil
}

// LoginValidator .
func LoginValidator() LoginAdmin {
	loginValidator := LoginAdmin{}
	return loginValidator
}
