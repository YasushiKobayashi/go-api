package user

import (
	"app/config"
	"app/handler"
	"app/models"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/scrypt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
)

func toHashFromScrypt(password string) string {
	salt := []byte(config.HASH_SALT)
	converted, _ := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
	return hex.EncodeToString(converted[:])
}

func Register() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		posts := new(models.User)
		if err = c.Bind(posts); err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}

		validate := validator.New()
		if err := validate.Struct(posts); err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}

		save := models.User{
			Id:       posts.Id,
			Name:     posts.Name,
			Email:    models.NewNullString(posts.Email.String),
			Password: models.NewNullString(posts.Password.String),
		}

		save, err = validation(save)
		if err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}

		data, err := models.CreateUser(save)
		if err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}
		strName := "username " + posts.Name
		strEmal := "emal " + posts.Email.String
		str := "が登録しました。"
		array := []string{strName, strEmal, str}
		toSlack := strings.Join(array, "\n")
		handler.SendSlack(toSlack)
		return c.JSON(http.StatusCreated, data)
	}
}

func validation(posts models.User) (res models.User, err error) {
	res = posts
	validate := validator.New()
	if err = validate.Struct(res); err != nil {
		return res, err
	}

	password := res.Password
	if !password.Valid || !res.Email.Valid {
		err = errors.New("password and email are required.")
		return res, err
	} else if utf8.RuneCountInString(password.String) < 8 {
		err = errors.New("password is short.")
		return res, err
	}

	res.Password.String = toHashFromScrypt(password.String)
	return res, err
}

func Login() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		param := new(models.User)
		if err = c.Bind(param); err != nil {
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}

		password := param.Password
		if password.Valid {
			password.String = toHashFromScrypt(password.String)
		}

		user := models.User{
			Email:    param.Email,
			Password: password,
		}
		data, err := models.Login(user)
		if err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}
		return c.JSON(http.StatusOK, data)
	}
}

func GetUserInfo(c echo.Context) (res models.UserJson, err error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)

	res = models.FindUser(claims.Id)
	if res.Id == 0 {
		return res, err
	}
	return res, err
}

func Get() echo.HandlerFunc {
	return func(c echo.Context) error {
		data, err := GetUserInfo(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}
		return c.JSON(http.StatusOK, data)
	}
}

func Update() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		posts := new(models.User)
		if err = c.Bind(posts); err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusNotAcceptable, config.NotAcceptable)
		}

		user, err := GetUserInfo(c)
		if err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}

		save := models.User{
			Id:       user.Id,
			Name:     posts.Name,
			Email:    models.NewNullString(posts.Email.String),
			Password: models.NewNullString(posts.Password.String),
			Image:    posts.Image,
			Created:  user.Created,
			Updated:  time.Now(),
		}

		save, err = validation(save)
		if err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}

		data, err := models.SaveUser(save)
		if err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}
		return c.JSON(http.StatusCreated, data)
	}
}

func Upload() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		file, err := c.FormFile("file")
		if err != nil {
			log.Printf("data : %v", err)
			return err
		}

		filePath, err := handler.Upload(file)
		if err != nil {
			log.Printf("data : %v", err)
			return err
		}

		user, err := GetUserInfo(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}

		image := models.NewNullString(filePath)
		param := models.User{
			Id:      user.Id,
			Name:    user.Name,
			Image:   image,
			Created: user.Created,
			Updated: time.Now(),
		}
		data, err := models.SaveUser(param)
		if err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}
		return c.JSON(http.StatusCreated, data)
	}
}
