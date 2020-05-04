package config

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/joho/godotenv"
)

var db *gorm.DB

func init() {
	var username string
	var password string
	var dbName string
	var dbHost string

	e := godotenv.Load()
	if e != nil {
		log.Print(e)
	}

	if os.Getenv("GO_ENV") != "production" {
		username = os.Getenv("DEV_DB_USER")
		password = os.Getenv("DEV_DB_PASS")
		dbName = os.Getenv("DEV_DB_NAME")
		dbHost = os.Getenv("DEV_DB_HOST")
	} else {
		username = os.Getenv("PROD_DB_USER")
		password = os.Getenv("PROD_DB_PASS")
		dbName = os.Getenv("PROD_DB_NAME")
		dbHost = os.Getenv("PROD_DB_HOST")
	}

	dbURI := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, dbHost, dbName)
	log.Println(dbURI)
	conn, err := gorm.Open("mysql", dbURI)

	if err != nil {
		log.Println(err)
		panic("failed to connect to database")
	}
	conn.LogMode(true)

	db = conn
}

// GetDB getdb
func GetDB() *gorm.DB {
	return db
}
