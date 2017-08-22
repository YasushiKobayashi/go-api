package models

import (
	"time"
)

type (
	Category struct {
		Id      int       `json:"id"`
		Name    string    `json:"name"`
		Slug    string    `json:"slug"`
		Posts   []Post    `json:"posts" gorm:"many2many:post_category;"`
		Created time.Time `json:"created" sql:"DEFAULT:current_timestamp"`
		Updated time.Time `json:"updated" sql:"DEFAULT:current_timestamp"`
	}

	CategoryJson struct {
		Id      int       `json:"id"`
		Name    string    `json:"name" validate:"required"`
		Slug    string    `json:"slug" validate:"required"`
		Created time.Time `json:"created" sql:"DEFAULT:current_timestamp"`
		Updated time.Time `json:"updated" sql:"DEFAULT:current_timestamp"`
	}

	PostCategory struct {
		PostID     int `gorm:"column:post_id"`
		CategoryID int `gorm:"column:category_id"`
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

func FindAllPostFromCategory(id int, pages int, search string) (res Category) {
	db.First(&res, id).
		Preload("User").Preload("Comments").
		Order("created desc").
		Limit(20).Offset(pages).
		Where("content LIKE ?", "%"+search+"%").
		Or("title LIKE ?", "%"+search+"%").
		Related(&res.Posts, "Posts")
	return res
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

func CountPostFromCategory(id int) (res Count) {
	var count int
	db.Model(&PostCategory{}).
		Where("category_id = ?", id).
		Count(&count)
	res.Count = count
	return res
}
