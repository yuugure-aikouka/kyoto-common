package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) healthCheckHandler(c echo.Context) error {
	return jsonResponse[any](c, http.StatusOK, nil)
}
