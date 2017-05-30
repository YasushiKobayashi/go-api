package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func init() {
	db := DB()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&CategoryJson{})
	db.AutoMigrate(&Comment{})
}
