package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Handler) HealthCheck(c echo.Context) error {
	message := "healthy"
	return jsonResponse(c, http.StatusOK, &message)
}
