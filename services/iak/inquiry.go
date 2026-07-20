package iak

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Inquiry(c echo.Context) error {
	// var (
	// 	svc     = "IAK Inquiry"
	// 	request models.RequestInquiry
	// )
	return c.JSON(http.StatusOK, nil)
}
