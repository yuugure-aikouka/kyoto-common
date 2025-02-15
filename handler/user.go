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
		return err
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
		return err
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
		return 0, jsonResponse[any](c, http.StatusBadRequest, nil)
	}

	_, err = h.store.GetUser(c.Request().Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, jsonResponse[any](c, http.StatusNotFound, nil, "User not found")
		}

		return 0, jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	return id, nil
}
