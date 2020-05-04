package transport

import (
	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"

	_adminService "uangbaik-account-microservice/admin/service"
)

// AdminRoute .
func AdminRoute(route *gin.RouterGroup, db *gorm.DB) {
	inDB := &_adminService.InDB{DB: db}
	v1 := route.Group("/admin")
	{
		v1.POST("/login", inDB.LoginHandler)
		v1.GET("/users/:id", inDB.GetUserDetail)
		v1.PATCH("/users/:id", inDB.EditUserData)
		v1.GET("/users", inDB.GetUserList)
		v1.POST("/bank", inDB.CreateBank)
		v1.GET("/bank", inDB.GetListBank)
		v1.GET("/bank/:id", inDB.GetDetailBank)
		v1.PATCH("/bank/:id", inDB.UpdateBank)
		v1.DELETE("/bank/:id", inDB.DeleteBank)
		v1.GET("/summary", inDB.Summary)
		v1.PATCH("/verify", inDB.VerifyStatus)
	}

	route.GET("/upload-asset", inDB.GetUploadAsset)
}
