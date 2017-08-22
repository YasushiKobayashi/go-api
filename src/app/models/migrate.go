package models

func init() {
	db := DB()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&CategoryJson{})
	db.AutoMigrate(&Comment{})
}
