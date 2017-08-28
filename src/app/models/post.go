package models

import (
	"log"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
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

	Search struct {
		Word string `json:"word" validate:"required"`
	}

	Upload struct {
		Path string `json:"path"`
	}
)

func FindAllPost() []Post {
	posts := []Post{}
	db.Preload("User").Preload("Comments.User").Order("created desc").Find(&posts)
	return posts
}

func FindPost(id int) Post {
	post := Post{}
	db.Preload("User").Preload("Comments.User").Preload("Categories").Find(&post, id)
	return post
}

func SearchPost(param Search) []Post {
	posts := []Post{}
	db.Preload("User").Preload("Comments.User").Order("created desc").Where("content LIKE ?", "%"+param.Word+"%").Or("title LIKE ?", "%"+param.Word+"%").Find(&posts)
	return posts
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
