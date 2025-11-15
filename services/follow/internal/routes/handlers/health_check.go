package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (f *FollowHandler) HealthCheck(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
}
