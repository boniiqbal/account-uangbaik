package transport

import (
	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"

	_userService "uangbaik-account-microservice/user/service"
)

// AuthenticationRoute .
func AuthenticationRoute(route *gin.RouterGroup, db *gorm.DB) {
	inDB := &_userService.InDB{DB: db}
	v1 := route.Group("/auth")
	{
		v1.POST("/register", inDB.CreateUser)
		v1.POST("/login", inDB.LoginUser)
		v1.POST("/login/refresh", inDB.RefreshToken)
		v1.POST("/login/verify", inDB.VerifyLogin)
		v1.POST("/pin", inDB.CreatePin)
		v1.POST("/pin/forgot", inDB.ForgotPin)
	}
}
 