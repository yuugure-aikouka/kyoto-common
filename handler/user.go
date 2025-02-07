package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetPartners(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return jsonResponse[any](c, http.StatusBadRequest, nil)
	}

	_, err = h.store.GetUser(c.Request().Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return jsonResponse[any](c, http.StatusNotFound, nil, "User not found")
		}

		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	partners, err := h.store.ListPartners(c.Request().Context(), int32(id))
	if err != nil {
		c.Logger().Errorf("Failed to fetch partners for user %d: %v", id, err)
		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	return jsonResponse(c, http.StatusOK, &partners)
}

func (h *Handler) GetPotentialPartners(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return jsonResponse[any](c, http.StatusBadRequest, nil)
	}

	_, err = h.store.GetUser(c.Request().Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return jsonResponse[any](c, http.StatusNotFound, nil, "User not found")
		}

		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	potentials, err := h.store.ListPotentialPartners(c.Request().Context(), int32(id))
	if err != nil {
		c.Logger().Errorf("Failed to fetch potential partners for user %d: %v", id, err)
		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	return jsonResponse(c, http.StatusOK, &potentials)
}
