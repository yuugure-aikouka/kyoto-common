package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/yuugure-aikouka/kyoto-common/utils"
)

func AuthMiddleware(tokenSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return utils.Response[any](c, http.StatusUnauthorized, nil, "Missing or invalid token")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := utils.ParseAuthToken(tokenString, tokenSecret)
			if err != nil {
				return utils.Response[any](c, http.StatusUnauthorized, nil, "Invalid or expired token")
			}

			claims := token.Claims.(jwt.MapClaims)
			c.Set("user_id", claims["user_id"])

			return next(c)
		}
	}
}
