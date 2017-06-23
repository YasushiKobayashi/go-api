package models

import (
	"app/config"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	validator "gopkg.in/go-playground/validator.v9"
)

type (
	User struct {
		Id       int        `db:"id" json:"id"`
		Name     string     `db:"name" json:"name" gorm:"not null" validate:"required"`
		Fbid     NullString `db:"fbid" json:"fbid"`
		Email    NullString `db:"email" json:"email" gorm:"unique_index" validate:"email"`
		Password NullString `db:"password" json:"password" validate:"min=8"`
		Image    NullString `db:"image" json:"image"`
		Created  time.Time  `json:"created" sql:"DEFAULT:current_timestamp"`
		Updated  time.Time  `json:"updated" sql:"DEFAULT:current_timestamp"`
	}

	UserJson struct {
		Id      int       `json:"id"`
		Name    string    `json:"name"`
		Email   *string   `json:"email"`
		Image   *string   `json:"image"`
		Created time.Time `json:"created"`
		Updated time.Time `json:"updated"`
	}

	JwtCustomClaims struct {
		Id int `json:"id"`
		jwt.StandardClaims
	}

	Token struct {
		Token string `json:"token"`
	}
)

func (UserJson) TableName() string {
	return "user"
}

func CreateUser(param User) (res Token, err error) {
	db := DB()
	user := param
	if err := db.Create(&user).Error; err != nil {
		log.Printf("data : %v", err)
		return res, err
	}

	token, err := createToken(user.Id)
	if err != nil {
		return res, err
	}

	res = Token{}
	res.Token = token
	return res, err
}

func Login(param User) (res Token, err error) {
	db := DB()
	user := param
	db.Where(&user).First(&user)

	if user.Id == 0 {
		return res, err
	}

	token, err := createToken(user.Id)
	if err != nil {
		return res, err
	}

	res = Token{}
	res.Token = token
	return res, err
}

func createToken(id int) (res string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = config.JWT_EXP
	res, err = token.SignedString([]byte("secret"))
	if err != nil {
		return res, err
	}
	return res, err
}

func SaveUser(params User) (res User, err error) {
	validate := validator.New()
	if err = validate.Struct(params); err != nil {
		log.Printf("data : %v", err)
		return res, err
	}

	db := DB()
	if err := db.Save(&params).Error; err != nil {
		return res, err
	}
	return params, err
}

func FindUser(id int) UserJson {
	db := DB()
	user := UserJson{}
	db.First(&user, id)
	return user
}
