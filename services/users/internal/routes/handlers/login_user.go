package handlers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kylehipz/socmed-microservices/users/internal/models"
	"github.com/kylehipz/socmed-microservices/users/internal/types"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserHandler) LoginUser(c echo.Context) error {
	var req types.LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid payload"})
	}

	var user models.User

	if err := u.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(u.jwtSecret))

	return c.JSON(http.StatusOK, echo.Map{
		"access_token": tokenString,
		"token_type":   "Bearer",
		"expires_in":   86400,
	})
}
