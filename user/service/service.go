package service

import (
	"github.com/jinzhu/gorm"
)

//InDB struct
type InDB struct {
	DB *gorm.DB
}
