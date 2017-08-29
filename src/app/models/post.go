package models

import (
	"log"
	"time"
)

type (
	Post struct {
		Id         int        `json:"id"`
		Title      string     `json:"title" gorm:"not null" validate:"required"`
		Content    string     `json:"content" gorm:"not null;type:TEXT" validate:"required"`
		WpFlg      bool       `json:"wp_flg" gorm:"not null"`
		UserId     int        `json:"user_id"`
		User       UserJson   `json:"user"`
		Comments   []Comment  `json:"comments"`
		Categories []Category `json:"categories" gorm:"many2many:post_category;"`
		Created    time.Time  `json:"created" sql:"DEFAULT:current_timestamp"`
		Updated    time.Time  `json:"updated" sql:"DEFAULT:current_timestamp"`
	}

	Upload struct {
		Path string `json:"path"`
	}
)

func SearchPost(pages int, search string) []Post {
	posts := []Post{}
	db.Limit(20).Offset(pages).Preload("User").Preload("Comments").
		Where("content LIKE ?", "%"+search+"%").
		Or("title LIKE ?", "%"+search+"%").
		Order("created desc").Find(&posts)
	return posts
}

func CountPost(search string) (res Count) {
	var count int
	db.Model(&Post{}).
		Where("content LIKE ?", "%"+search+"%").Or("title LIKE ?", "%"+search+"%").
		Count(&count)
	res.Count = count
	return res
}

func FindPost(id int) Post {
	post := Post{}
	db.Preload("User").Preload("Comments.User").Preload("Categories").Find(&post, id)
	return post
}

func UsersPost(id int) []Post {
	posts := []Post{}
	db.Preload("User").Preload("Comments.User").Order("created desc").Where(Post{UserId: id}).Find(&posts)
	return posts
}

func CreatePost(params Post) (res Post, err error) {
	if err = validate.Struct(params); err != nil {
		log.Printf("data : %v", err)
		return res, err
	}

	if err := db.Create(&params).Related(&params.Categories, "Categories").Error; err != nil {
		return res, err
	}
	return params, err
}

func SavePost(params Post) (res Post, err error) {
	if err = validate.Struct(params); err != nil {
		log.Printf("data : %v", err)
		return res, err
	}

	if err := db.Save(&params).Related(&params.Categories, "Categories").Error; err != nil {
		return res, err
	}
	return params, err
}
