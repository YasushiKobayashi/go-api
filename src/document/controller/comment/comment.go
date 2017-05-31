package comment

import (
	"document/config"
	"document/handler"
	"document/models"
	"log"
	"net/http"
	"strconv"
	"strings"
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

		toSlack := config.SITE_TITLE + "ドキュメントの更新"
		toSlack, err = slackSendCont(toSlack, data)
		if err != nil {
			log.Printf("data : %v", "SLACKの通知エラー")
			log.Printf("data : %v", err)
		}
		handler.SendSlack(toSlack)
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
			Id:      int(id),
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

		toSlack := config.SITE_TITLE + "ドキュメントの更新"
		toSlack, err = slackSendCont(toSlack, data)
		if err != nil {
			log.Printf("data : %v", "SLACKの通知エラー")
			log.Printf("data : %v", err)
		}
		handler.SendSlack(toSlack)
		return c.JSON(http.StatusOK, data)
	}
}

func slackSendCont(str string, param models.Comment) (res string, err error) {
	user := models.FindUser(param.UserId)
	if user.Id == 0 {
		return res, err
	}

	userStr := "_by " + user.Name + "_"
	urlStr := "url " + config.FRONT_URL + "article/" + strconv.Itoa(param.PostId)
	contStr := handler.TrimStr(param.Content, 40)
	array := []string{str, userStr, urlStr, contStr}
	res = strings.Join(array, "\n")
	return res, err
}
