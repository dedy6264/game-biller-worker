package main

import (
	"fmt"
	"game-biller-worker/configs"
	"game-biller-worker/routes"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	fmt.Println("====", configs.APP_V)
	// Initialize DB Connection
	// ss := helpers.IAKConverterResponse(helpers.ProviderPayload{
	// 	ResponseCode: "06",
	// 	Message:      "INQUIRY ID NOT FOUND",
	// 	Description:  "The inquiry ID (tr_id) that youve inputted is not found, there is no inquiry with that ID. You can check the inquiry ID (tr_id) field for any typos, or try using another inquiry ID.",
	// 	Case:         "PAY",
	// }, 200)
	// s, _ := json.Marshal(ss)
	// fmt.Println(string(s))
	e := echo.New()

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	// Register Routes
	routes.AppRoutes(e)

	// Start Server
	e.Logger.Fatal(e.Start(":" + configs.APP_PORT))
}
