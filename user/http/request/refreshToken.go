package http

import (
	_config "uangbaik-account-microservice/config"

	"github.com/gin-gonic/gin"
)

//RefreshToken payload struct
type RefreshToken struct {
	Data struct {
		RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
	} `json:"data"`
}

//BindRefreshToken for bind request
func (refresh *RefreshToken) BindRefreshToken(c *gin.Context) error {
	err := _config.BindRequest(c, refresh)
	if err != nil {
		return err
	}
	return nil
}

// RefreshTokenValidator .
func RefreshTokenValidator() RefreshToken {
	refreshTokenValidator := RefreshToken{}
	return refreshTokenValidator
}
