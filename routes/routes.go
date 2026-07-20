package routes

import (
	"bytes"
	"encoding/json"
	"game-biller-worker/services/iak"

	"io"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func AppRoutes(e echo.Echo) {
	// Print logs for every request, from header, request and response
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// --- 1. Log Request Headers ---
			reqHeaderBytes, _ := json.Marshal(c.Request().Header)
			log.Println("Request Endpoint :: ", c.Request().URL.Path)
			log.Println("Request Headers :: ", string(reqHeaderBytes))

			// --- 2. Safely Log Request Body ---
			var reqBodyBytes []byte
			if c.Request().Body != nil {
				// Read the body bytes
				reqBodyBytes, _ = io.ReadAll(c.Request().Body)
				// IMPORTANT: Restore the body so downstream handlers can read it too!
				c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
			}
			log.Println("Request Body :: ", string(reqBodyBytes))

			// --- 3. Set up Response Interceptor ---
			// We wrap the response writer to capture what goes out
			// --- 3. Set up Response Interceptor ---
			resBodyBuffer := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBodyBuffer)

			// Explicitly name the fields so the compiler maps them correctly
			writer := &responseBodyWriter{
				ResponseWriter: c.Response().Writer, // The original http.ResponseWriter
				Writer:         mw,                  // The multi-writer stream
			}
			c.Response().Writer = writer

			// Defer the response logging until after next(c) executes
			defer func() {
				resHeaderBytes, _ := json.Marshal(c.Response().Header())
				log.Println("Response Headers :: ", string(resHeaderBytes))
				log.Println("Response Body :: ", resBodyBuffer.String())
			}()

			return next(c)
		}
	})

	Iak(e)
}
func Iak(e echo.Echo) {
	Iak := e.Group("/api/iak")
	Iak.POST("/inquiry", iak.Inquiry)
	Iak.POST("/payment", iak.Payment)
}
func timeNowStr() string {
	return time.Now().Format(time.RFC3339)
}

type responseBodyWriter struct {
	http.ResponseWriter
	io.Writer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
