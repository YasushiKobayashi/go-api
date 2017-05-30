package models

import (
	"log"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	validator "gopkg.in/go-playground/validator.v9"
)

type (
	Post struct {
		Id         uint       `json:"id"`
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
	db := DB()
	posts := []Post{}
	db.Find(&posts)
	for i, _ := range posts {
		db.Model(posts[i]).Related(&posts[i].User, "User").Related(&posts[i].Comments, "Comments")
	}
	return posts
}

func FindPost(id int) Post {
	db := DB()
	post := Post{}
	db.First(&post, id).Related(&post.User, "User").Related(&post.Categories, "Categories").Related(&post.Comments)

	for i, _ := range post.Comments {
		db.Model(post.Comments[i]).Related(&post.Comments[i].User, "User")
	}
	return post
}

func SearchPost(param Search) []Post {
	db := DB()
	posts := []Post{}
	db.Where("content LIKE ?", "%"+param.Word+"%").Or("title LIKE ?", "%"+param.Word+"%").Find(&posts)
	for i, _ := range posts {
		db.Model(posts[i]).Related(&posts[i].User, "User").Related(&posts[i].Comments, "Comments")
	}
	return posts
}

func UsersPost(id int) []Post {
	db := DB()
	posts := []Post{}
	db.Where(Post{UserId: id}).Find(&posts)
	for i, _ := range posts {
		db.Model(posts[i]).Related(&posts[i].User, "User").Related(&posts[i].Comments, "Comments")
	}
	return posts
}

func CreatePost(params Post) (res Post, err error) {
	validate := validator.New()
	if err = validate.Struct(params); err != nil {
		log.Printf("data : %v", err)
		return res, err
	}

	db := DB()
	if err := db.Create(&params).Related(&params.Categories, "Categories").Error; err != nil {
		return res, err
	}
	return params, err
}

func SavePost(params Post) (res Post, err error) {
	validate := validator.New()
	if err = validate.Struct(params); err != nil {
		log.Printf("data : %v", err)
		return res, err
	}

	db := DB()
	if err := db.Save(&params).Related(&params.Categories, "Categories").Error; err != nil {
		return res, err
	}
	return params, err
}
