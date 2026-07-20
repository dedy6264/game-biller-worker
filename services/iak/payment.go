package iak

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Payment(c echo.Context) error {
	// var (
	// 	svc     = "IAK Payment"
	// 	request models.RequestInquiry
	// )
	return c.JSON(http.StatusOK, nil)
}
