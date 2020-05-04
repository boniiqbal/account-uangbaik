package config

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//BindRequest for
func BindRequest(c *gin.Context,  obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(obj, b)
}