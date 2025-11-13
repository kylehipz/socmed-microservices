package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (u *UserHandler) HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
