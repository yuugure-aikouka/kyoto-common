package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetPartners(c echo.Context) error {
	id, err := h.getIdAndValidateUser(c)
	if err != nil {
		httpErr := err.(*echo.HTTPError)
		return jsonResponse[any](c, httpErr.Code, nil, httpErr.Message)
	}

	partners, err := h.store.ListPartners(c.Request().Context(), int32(id))
	if err != nil {
		c.Logger().Errorf("Failed to fetch partners for user %d: %v", id, err)
		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	return jsonResponse(c, http.StatusOK, &partners)
}

func (h *Handler) GetPotentialPartners(c echo.Context) error {
	id, err := h.getIdAndValidateUser(c)
	if err != nil {
		httpErr := err.(*echo.HTTPError)
		return jsonResponse[any](c, httpErr.Code, nil, httpErr.Message)
	}

	potentials, err := h.store.ListPotentialPartners(c.Request().Context(), int32(id))
	if err != nil {
		c.Logger().Errorf("Failed to fetch potential partners for user %d: %v", id, err)
		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	return jsonResponse(c, http.StatusOK, &potentials)
}

func (h *Handler) getIdAndValidateUser(c echo.Context) (int, error) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest)
	}

	_, err = h.store.GetUserById(c.Request().Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, echo.NewHTTPError(http.StatusNotFound, "User not found")
		}

		return 0, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return id, nil
}
