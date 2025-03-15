package handler

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	db "github.com/yuugure-aikouka/kyoto-common/db/store"
	"github.com/yuugure-aikouka/kyoto-common/model"
	"github.com/yuugure-aikouka/kyoto-common/utils"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) SendMagicLinkHandler(c echo.Context) error {
	var request model.SendMagicLinkRequest
	if err := c.Bind(&request); err != nil {
		return jsonResponse[any](c, http.StatusBadRequest, nil, "Need \"email\"")
	}

	if !isValidEmail(request.Email) {
		return jsonResponse[any](c, http.StatusBadRequest, nil, "Invalid email")
	}

	ctx := c.Request().Context()

	if err := h.store.InvalidateActiveMagicLink(ctx, request.Email); err != nil {
		c.Logger().Errorf("Failed to invalidate magic link %s: %v", request.Email, err)
		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	rawToken := h.createMagicLinkToken()
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(rawToken), bcrypt.DefaultCost)
	if err != nil {
		c.Logger().Errorf("Failed to hash magic link token %s: %v", request.Email, err)
		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	createMagicLinkParams := db.CreateMagicLinkParams{Email: request.Email, Token: string(hashedToken)}
	if err := h.store.CreateMagicLink(ctx, createMagicLinkParams); err != nil {
		c.Logger().Errorf("Failed to create magic link record %s: %v", request.Email, err)
		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	magicLink := h.createMagicLink(rawToken)
	if err := h.mailer.SendMagicLinkMail(request.Email, magicLink); err != nil {
		c.Logger().Errorf("Failed sending magic link to %s: %v", request.Email, err)
		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	message := fmt.Sprintf("Magic link has been sent to %s", request.Email)
	return jsonResponse(c, http.StatusOK, &message)
}

func (h *Handler) createMagicLinkToken() string {
	return utils.RandomBase64EncodedString(32) // 32 bytes -> 43-character token
}

func (h *Handler) createMagicLink(token string) string {
	appUrl := h.config.AppBaseUrl

	return fmt.Sprintf("%s/auth/magic?token=%s", appUrl, token)
}

// TODO: use the validator from joshua later
func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
