package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIResponse struct {
	Status string        `json:"status"`
	Data   interface{}   `json:"data,omitempty"`
	Errors []interface{} `json:"errors,omitempty"`
}

func jsonResponse(c echo.Context, code int, data interface{}, errors ...interface{}) error {
	return c.JSON(code, APIResponse{
		Status: http.StatusText(code),
		Data:   data,
		Errors: errors,
	})
}
