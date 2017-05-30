package comment

import (
	"document/config"
	"document/models"
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func Create() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		params := new(models.Comment)
		if err = c.Bind(params); err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusNotAcceptable, config.NotAcceptable)
		}

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*models.JwtCustomClaims)

		comment := models.Comment{
			UserId:  claims.Id,
			PostId:  params.PostId,
			Content: params.Content,
		}
		data, err := models.CreateComment(comment)
		if err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, config.BadRequest)
		}
		return c.JSON(http.StatusCreated, data)
	}
}

func Update() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		params := new(models.Comment)
		log.Printf("data : %v", params)
		if err = c.Bind(params); err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusNotAcceptable, config.NotAcceptable)
		}

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*models.JwtCustomClaims)

		comment := models.FindComment(int(id))
		if comment.UserId != claims.Id {
			return c.JSON(http.StatusUnauthorized, config.Unauthorized)
		}

		post := models.Comment{
			Id:      uint(id),
			UserId:  claims.Id,
			PostId:  comment.PostId,
			Content: params.Content,
			Created: comment.Created,
			Updated: time.Now(),
		}
		data, err := models.SaveComment(post)
		if err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, data)
	}
}
