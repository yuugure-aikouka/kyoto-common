package api

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type GetPartnersRequest struct {
	ID int32 `param:"id" validate:"required,min=1"`
}

func (s *Server) getPartnersHandler(c echo.Context) error {
	var req GetPartnersRequest
	if err := ValidateRequest(c, &req); err != nil {
		return jsonResponse(c, http.StatusBadRequest, nil, err.Error())
	}

	_, err := s.store.GetUser(c.Request().Context(), req.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return jsonResponse(c, http.StatusNotFound, nil, "User not found")
		}

		return jsonResponse(c, http.StatusInternalServerError, nil)
	}

	partners, err := s.store.ListPartners(c.Request().Context(), req.ID)
	if err != nil {
		c.Logger().Errorf("Failed to fetch partners for user %d: %v", req.ID, err)
		return jsonResponse(c, http.StatusInternalServerError, nil)
	}

	return jsonResponse(c, http.StatusOK, partners)
}
