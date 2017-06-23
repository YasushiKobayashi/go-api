package category

import (
	"app/config"
	"app/models"
	"net/http"

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
