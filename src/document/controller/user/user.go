package user

import (
	"document/config"
	"document/handler"
	"document/models"
	"encoding/hex"
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
		param := new(models.User)
		if err = c.Bind(param); err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}

		validate := validator.New()
		if err := validate.Struct(param); err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}

		password := param.Password
		if !password.Valid || !param.Email.Valid {
			log.Printf("data : %v", "password and email are required.")
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		} else if utf8.RuneCountInString(password.String) < 8 {
			log.Printf("data : %v", "password is short.")
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}

		if password.Valid {
			password.String = toHashFromScrypt(password.String)
		}

		user := models.User{
			Id:       param.Id,
			Name:     param.Name,
			Fbid:     param.Fbid,
			Email:    param.Email,
			Password: password,
		}

		data, err := models.CreateUser(user)
		if err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}
		strName := "username " + param.Name
		strEmal := "emal " + param.Email.String
		str := "が登録しました。"
		array := []string{strName, strEmal, str}
		toSlack := strings.Join(array, "\n")
		handler.SendSlack(toSlack)
		return c.JSON(http.StatusCreated, data)
	}
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

func Get() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*models.JwtCustomClaims)

		data := models.FindUser(claims.Id)
		if data.Id == 0 {
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}
		return c.JSON(http.StatusOK, data)
	}
}

func Put() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusBadRequest, config.BadRequest)
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
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*models.JwtCustomClaims)

		userInfo := models.FindUser(claims.Id)
		if userInfo.Id == 0 {
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}

		image := models.NewNullString(filePath)
		param := models.User{
			Id:      int(claims.Id),
			Name:    userInfo.Name,
			Image:   image,
			Created: userInfo.Created,
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
