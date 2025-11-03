package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "http://keycloak:8080/realms/socmed-microservices")
	if err != nil {
		panic(err)
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: "kong-client", // optional: verify audience
		SkipIssuerCheck: true, // remove in prod
		SkipClientIDCheck: true,  // remove in prod
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing bearer token")
			}

			rawToken := strings.TrimPrefix(auth, "Bearer ")

			// Verify and decode token
			idToken, err := verifier.Verify(ctx, rawToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token: "+err.Error())
			}

			// Extract claims
			var claims map[string]any
			if err := idToken.Claims(&claims); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid claims")
			}

			c.Set("user", claims)
			return next(c)
		}
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "HELLO WORLD!!")
	})

	e.GET("/me", func(c echo.Context) error {
		return c.JSON(http.StatusOK, c.Get("user"))
	})

	e.Logger.Fatal(e.Start(":8080"))
}
