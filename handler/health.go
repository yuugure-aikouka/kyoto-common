package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yuugure-aikouka/kyoto-common/utils"
)

func (s *Handler) HealthCheck(c echo.Context) error {
	message := "healthy"
	return utils.Response(c, http.StatusOK, &message)
}
