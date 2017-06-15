package main

import (
	"document/config"
	"document/controller/category"
	"document/controller/comment"
	"document/controller/post"
	"document/controller/user"
	"document/models"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static", "static")
	e.Use(middleware.CORS())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{config.ALLOW_ORIGINS},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.POST("/v1/register", user.Register())
	e.POST("/v1/login", user.Login())

	r := e.Group("/v1/")
	jwtconfig := middleware.JWTConfig{
		Claims:     &models.JwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	r.Use(middleware.JWTWithConfig(jwtconfig))
	r.GET("user", user.Get())
	r.PUT("user", user.Put())
	r.POST("user/upload", user.Upload())

	r.GET("post", post.List())
	r.POST("post", post.Create())
	r.GET("post/:id", post.Get())
	r.PUT("post/:id", post.Update())
	r.POST("post/search", post.Search())
	r.GET("post/user", post.GetFromUser())
	r.GET("post/category/:category_id", post.GetFromCategory())
	r.POST("post/upload", post.Upload())

	r.POST("comment", comment.Create())
	r.PUT("comment/:id", comment.Update())

	r.GET("category", category.List())
	r.POST("category", category.Create())

	log.Fatal(e.Start(config.HOST + config.PORT))
}
