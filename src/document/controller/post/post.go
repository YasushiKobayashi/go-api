package post

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

func List() echo.HandlerFunc {
	return func(c echo.Context) error {
		data := models.FindAllPost()
		return c.JSON(http.StatusOK, data)
	}
}

func Get() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		data := models.FindPost(int(id))
		if data.Id == 0 {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusNotFound, config.NotFound)
		}
		return c.JSON(http.StatusOK, data)
	}
}

func Search() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		params := new(models.Search)
		if err = c.Bind(params); err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusNotAcceptable, config.NotAcceptable)
		}

		search := models.Search{
			Word: params.Word,
		}
		data := models.SearchPost(search)
		return c.JSON(http.StatusOK, data)
	}
}

func GetFromCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("category_id"))
		if err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		data := models.FindAllPostFromCategory(int(id))
		return c.JSON(http.StatusOK, data.Posts)
	}
}

func GetFromUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*models.JwtCustomClaims)

		data := models.UsersPost(claims.Id)
		return c.JSON(http.StatusOK, data)
	}
}

func Create() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		params := new(models.Post)
		if err = c.Bind(params); err != nil {
			return c.JSON(http.StatusNotAcceptable, config.NotAcceptable)
		}

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*models.JwtCustomClaims)

		post := models.Post{
			UserId:     claims.Id,
			Title:      params.Title,
			Content:    params.Content,
			WpFlg:      params.WpFlg,
			Categories: params.Categories,
		}
		data, err := models.CreatePost(post)
		if err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		toSlack := config.SITE_TITLE + "新規ドキュメント"
		toSlack, err = slackSendCont(toSlack, claims.Id, data)
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

		params := new(models.Post)
		if err = c.Bind(params); err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusNotAcceptable, config.NotAcceptable)
		}

		nowPost := models.FindPost(int(id))
		if nowPost.Id == 0 {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusNotFound, config.NotFound)
		}

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*models.JwtCustomClaims)

		post := models.Post{
			Id:         int(id),
			UserId:     claims.Id,
			Title:      params.Title,
			Content:    params.Content,
			WpFlg:      params.WpFlg,
			Categories: params.Categories,
			Created:    nowPost.Created,
			Updated:    time.Now(),
		}
		data, err := models.SavePost(post)
		if err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		toSlack := config.SITE_TITLE + "ドキュメントの更新"
		toSlack, err = slackSendCont(toSlack, claims.Id, data)
		if err != nil {
			log.Printf("data : %v", "SLACKの通知エラー")
			log.Printf("data : %v", err)
		}
		handler.SendSlack(toSlack)
		return c.JSON(http.StatusOK, data)
	}
}

func slackSendCont(str string, id int, param models.Post) (res string, err error) {
	user := models.FindUser(id)
	if user.Id == 0 {
		return res, err
	}

	ttlStr := param.Title
	userStr := "_by " + user.Name + "_"
	urlStr := "url " + config.FRONT_URL + "article/" + strconv.Itoa(param.Id)
	contStr := handler.TrimStr(param.Content, 40)
	array := []string{str, ttlStr, userStr, urlStr, contStr}
	res = strings.Join(array, "\n")
	return res, err
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

		data := models.Upload{
			Path: filePath,
		}
		return c.JSON(http.StatusCreated, data)
	}
}
