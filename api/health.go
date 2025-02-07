package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) healthCheckHandler(c echo.Context) error {
	message := "healthy"
	return jsonResponse(c, http.StatusOK, &message)
}
