package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(echoContext echo.Context) error {
		authorizationHeader := echoContext.Request().Header.Get("Authorization")

		if authorizationHeader == "" {
			echoContext.String(http.StatusNetworkAuthenticationRequired, "Authorization header is required")

			return nil
		}
		if !ApiKeyStore.IsValid(authorizationHeader) {
			echoContext.String(http.StatusUnauthorized, "Unauthorized request")

			return nil
		}

		next(echoContext)
		return nil
	}
}
