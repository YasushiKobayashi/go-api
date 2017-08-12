package category

import (
	"app/config"
	"app/models"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func List() echo.HandlerFunc {
	return func(c echo.Context) error {
		data := models.FindAllCategory()
		return c.JSON(http.StatusOK, data)
	}
}

func Create() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		params := new(models.CategoryJson)
		if err = c.Bind(params); err != nil {
			return c.JSON(http.StatusNotAcceptable, config.NotAcceptable)
		}

		category := models.CategoryJson{
			Slug: params.Slug,
			Name: params.Name,
		}
		data, err := models.CreateCategory(category)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusCreated, data)
	}
}

func GetWithPostList() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("category_id"))
		if err != nil {
			log.Printf("data : %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		query := c.QueryParam("q")
		pages := c.QueryParam("pages")
		if pages == "" {
			pages = "1"
		}

		var number int
		number, _ = strconv.Atoi(pages)
		log.Printf("data : %v", pages)
		data := models.FindAllPostFromCategory(int(id), (int(number)-1)*20, query)
		return c.JSON(http.StatusOK, data)
	}
}
