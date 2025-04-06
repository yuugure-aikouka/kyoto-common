package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	db "github.com/yuugure-aikouka/kyoto-common/db/store"
	"github.com/yuugure-aikouka/kyoto-common/model"
	"github.com/yuugure-aikouka/kyoto-common/utils"
)

const (
	MagicLinkTokenLength         = 32
	MagicLinkCacheKeyPrefix      = "auth:magic_link"
	MagicLinkCacheExpiryDuration = 15 * time.Minute
	AccessTokenDuration          = 15 * time.Minute
	RefreshTokenDuration         = 30 * (24 * time.Hour)
)

func (h *Handler) SendMagicLinkHandler(c echo.Context) error {
	ctx := c.Request().Context()

	var request model.SendMagicLinkRequest
	if err := c.Bind(&request); err != nil || request.Email == "" { // TODO: use validator from joshua
		return jsonResponse[any](c, http.StatusBadRequest, nil, "Missing \"email\" field")
	}

	if !isValidEmail(request.Email) {
		return jsonResponse[any](c, http.StatusBadRequest, nil, "Invalid email")
	}

	token := h.createMagicLinkToken()
	cacheKey := h.getTokenCacheKey(token)
	if err := h.cache.Set(ctx, cacheKey, request.Email, MagicLinkCacheExpiryDuration); err != nil {
		c.Logger().Errorf("Failed to create magic link record %s: %v", request.Email, err)
		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	magicLink := h.createMagicLink(token)
	if err := h.mailer.SendMagicLinkMail(request.Email, magicLink); err != nil {
		c.Logger().Errorf("Failed sending magic link to %s: %v", request.Email, err)
		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	message := fmt.Sprintf("Magic link has been sent to %s", request.Email)
	return jsonResponse(c, http.StatusOK, &message)
}

func (h *Handler) VerifyMagicLinkHandler(c echo.Context) error {
	ctx := c.Request().Context()
	token := c.QueryParam("token")
	if token == "" {
		return jsonResponse[any](c, http.StatusBadRequest, nil, "Need \"token\" query parameter")
	}

	cacheKey := h.getTokenCacheKey(token)
	userEmail, err := h.cache.Get(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return jsonResponse[any](c, http.StatusUnauthorized, nil, "Invalid token or expired")
		}

		c.Logger().Errorf("Cache retrieval error for key: %s: %v", cacheKey, err)
		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	user, err := h.store.GetUserByEmail(ctx, userEmail)
	if err != nil {
		isNewUser := errors.Is(err, pgx.ErrNoRows)
		if isNewUser {
			user, err = h.store.CreateUser(ctx, db.CreateUserParams{
				Email:       userEmail,
				Username:    "",
				DisplayName: "",
				AvatarUrl:   "",
				IsAi:        false,
			})
			if err != nil {
				c.Logger().Errorf("Failed to create user %s: %v", userEmail, err)
				return jsonResponse[any](c, http.StatusInternalServerError, nil)
			}
		} else {
			c.Logger().Errorf("Failed checking user %s: %v", userEmail, err)
			return jsonResponse[any](c, http.StatusInternalServerError, nil)
		}
	}

	accessToken, err := h.createAuthToken(user.ID, AccessTokenDuration, h.config.AccessTokenSecretKey)
	if err != nil {
		c.Logger().Errorf("Failed to create access token for user %d: %v", user.ID, err)
		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	refreshToken, err := h.createAuthToken(user.ID, RefreshTokenDuration, h.config.RefreshTokenSecretKey)
	if err != nil {
		c.Logger().Errorf("Failed to create refresh token for user %d: %v", user.ID, err)
		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	err = h.store.CreateRefreshToken(ctx, refreshToken)
	if err != nil {
		c.Logger().Errorf("Failed to save refresh token %s: %v", user.ID, err)
		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	h.cache.Delete(ctx, cacheKey)
	return jsonResponse(c, http.StatusOK, &model.VerifyMagicLinkResponseBody{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *Handler) RefreshAccessTokenHandler(c echo.Context) error {
	ctx := c.Request().Context()

	var request model.RefreshAccessTokenRequest
	if err := c.Bind(&request); err != nil || request.RefreshToken == "" {
		return jsonResponse[any](c, http.StatusBadRequest, nil, "Missing \"refresh_token\" field")
	}

	refreshToken, err := h.parseAuthToken(request.RefreshToken, h.config.RefreshTokenSecretKey)
	if err != nil || !refreshToken.Valid {
		return jsonResponse[any](c, http.StatusUnauthorized, nil, "Invalid or expired token")
	}

	exist, err := h.store.CheckRefreshTokenExist(ctx, refreshToken.Raw)
	if err != nil {
		c.Logger().Errorf("Failed to check refresh token %s: %v", refreshToken.Raw, err)
		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	if !exist {
		return jsonResponse[any](c, http.StatusUnauthorized, nil, "Invalid or expired token")
	}

	claims := refreshToken.Claims.(jwt.MapClaims)
	userId := int32(claims["user_id"].(float64))

	accessToken, err := h.createAuthToken(userId, AccessTokenDuration, h.config.AccessTokenSecretKey)
	if err != nil {
		c.Logger().Errorf("Failed to create access token for user %d: %v", userId, err)
		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	return jsonResponse(c, http.StatusOK, &model.RefreshAccessTokenResponseBody{AccessToken: accessToken})
}

func (h *Handler) LogoutHandler(c echo.Context) error {
	var request model.LogoutRequest
	if err := c.Bind(&request); err != nil || request.RefreshToken == "" {
		return jsonResponse[any](c, http.StatusBadRequest, nil, "Missing \"refresh_token\" field")
	}

	err := h.store.DeleteRefreshToken(context.Background(), request.RefreshToken)
	if err != nil {
		c.Logger().Errorf("Failed to delete refresh token %s: %v", request.RefreshToken, err)
		return jsonResponse[any](c, http.StatusInternalServerError, nil)
	}

	message := "Logout successful"
	return jsonResponse(c, http.StatusOK, &message)
}

func (h *Handler) createMagicLinkToken() string {
	return utils.RandomBase64EncodedString(MagicLinkTokenLength)
}

func (h *Handler) createMagicLink(token string) string {
	appUrl := h.config.AppBaseUrl

	return fmt.Sprintf("%s/auth/magic?token=%s", appUrl, token)
}

func (h *Handler) getTokenCacheKey(token string) string {
	hashInBytes := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hashInBytes[:])

	return fmt.Sprintf("%s:%s", MagicLinkCacheKeyPrefix, hashedToken)
}

func (h *Handler) createAuthToken(userId int32, expiryDuration time.Duration, key string) (string, error) {
	now := time.Now()
	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":     h.config.AppBaseUrl,
		"sub":     userId,
		"exp":     now.Add(expiryDuration).Unix(),
		"iat":     now.Unix(),
		"jti":     uuid.New().String(),
		"user_id": userId,
	})

	return authToken.SignedString([]byte(key))
}

func (h *Handler) parseAuthToken(tokenString, key string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
}

// TODO: use the validator from joshua later
func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
