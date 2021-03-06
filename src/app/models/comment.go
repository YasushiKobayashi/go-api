package models

import (
	"log"
	"time"
)

type (
	Comment struct {
		Id      int       `json:"id"`
		UserId  int       `json:"user_id" gorm:"not null" validate:"required"`
		User    UserJson  `json:"user"`
		PostId  int       `json:"post_id" gorm:"not null"`
		Content string    `json:"content" gorm:"not null;type:TEXT" validate:"required"`
		Created time.Time `json:"created" sql:"DEFAULT:current_timestamp"`
		Updated time.Time `json:"updated" sql:"DEFAULT:current_timestamp"`
	}
)

func CreateComment(params Comment) (res Comment, err error) {
	if err = validate.Struct(params); err != nil {
		log.Printf("data : %v", err)
		return res, err
	}

	if err := db.Create(&params).Error; err != nil {
		log.Printf("data : %v", err)
		return res, err
	}
	if err := db.Model(params).Related(&params.User, "User").Error; err != nil {
		log.Printf("data : %v", err)
		return res, err
	}
	return params, err
}

func SaveComment(params Comment) (res Comment, err error) {
	if err = validate.Struct(params); err != nil {
		log.Printf("data : %v", err)
		return res, err
	}

	if err := db.Save(&params).Error; err != nil {
		log.Printf("data : %v", err)
		return res, err
	}
	return params, err
}

func FindComment(id int) Comment {
	comment := Comment{}
	db.First(&comment, id)
	return comment
}
