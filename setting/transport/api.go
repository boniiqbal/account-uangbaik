package transport

import (
	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"

	_settingService "uangbaik-account-microservice/setting/service"
)

// SettingRoute .
func SettingRoute(route *gin.RouterGroup, db *gorm.DB) {
	inDB := &_settingService.InDB{DB: db}
	v1 := route.Group("/setting")
	{
		v1.GET("/me", inDB.GetMe)
		v1.PATCH("/", inDB.EditProfile)
		v1.POST("/upload", inDB.Upload)
		v1.GET("/bank", inDB.GetListBank)
		v1.POST("/bank-account", inDB.CreateBankAccount)
		v1.GET("/bank-account", inDB.GetListBankAccount)
		v1.GET("/bank-account/:id", inDB.GetDetailBankAccount)
		v1.DELETE("/bank-account/:id", inDB.DeleteBankAccount)
	}
}
