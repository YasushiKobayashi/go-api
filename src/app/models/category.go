package models

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type (
	Category struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Slug  string `json:"slug"`
		Posts []Post `json:"posts" gorm:"many2many:post_category;"`
	}

	CategoryJson struct {
		Id      int       `json:"id"`
		Name    string    `json:"name" validate:"required"`
		Slug    string    `json:"slug" validate:"required"`
		Created time.Time `json:"created" sql:"DEFAULT:current_timestamp"`
		Updated time.Time `json:"updated" sql:"DEFAULT:current_timestamp"`
	}
)

func (CategoryJson) TableName() string {
	return "category"
}

func FindAllCategory() []Category {
	category := []Category{}
	db.Find(&category)
	return category
}

func FindAllPostFromCategory(id int) Category {
	category := Category{}
	db.First(&category, id).Order("created desc").Preload("User").Preload("Comments.User").Related(&category.Posts, "Posts")
	return category
}

func CreateCategory(params CategoryJson) (res CategoryJson, err error) {
	if err = validate.Struct(params); err != nil {
		return res, err
	}

	if err := db.Create(&params).Error; err != nil {
		return res, err
	}
	return params, err
}
