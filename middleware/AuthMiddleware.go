package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/yasngleer/studentplan/utils"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get  Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is missing")
			}

			// Check if the header starts with "Bearer "
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token format")
			}

			tokenString := parts[1]

			// Parse the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(utils.Jwtsecret), nil
			})

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID := claims["id"].(string)
				c.Set("user_id", userID)
				return next(c)
			}

			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}
	}
}
