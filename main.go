package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_admin "uangbaik-account-microservice/admin/transport"
	"uangbaik-account-microservice/config"
	_setting "uangbaik-account-microservice/setting/transport"
	_user "uangbaik-account-microservice/user/transport"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	db := config.GetDB()
	v1 := router.Group("/api/v1")

	_user.AuthenticationRoute(v1, db)
	_admin.AdminRoute(v1, db)
	_setting.SettingRoute(v1, db)

	// Health check
	v1.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "I'm well",
		})
	})

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal(fmt.Sprintf("PORT must be set [%s]", port))
	}

	router.Run(":" + port)
}
