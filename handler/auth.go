package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prokobit/auth-service/model"
	account "github.com/prokobit/auth-service/store"
)

type AuthHandler struct {
	store account.Store
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{account.NewMemoryStore()}
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func NewAuthResponse(token string) *AuthResponse {
	return &AuthResponse{
		Token: token,
	}
}

func (h *AuthHandler) SignUp(c echo.Context) error {
	r := new(SignUpRequest)
	if err := c.Bind(r); err != nil {
		return err
	}
	a, err := model.NewAccount(r.Email, r.Password)
	if err != nil {
		return err
	}
	if err := h.store.Create(a); err != nil {
		return nil
	}
	t, err := a.GenerateJWT()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, NewAuthResponse(t))
}
func (h *AuthHandler) Login(c echo.Context) error {
	r := new(LoginRequest)
	if err := c.Bind(r); err != nil {
		return err
	}
	a, err := h.store.GetByEmail(r.Email)
	if err != nil {
		return echo.ErrUnauthorized
	}
	if !a.CheckPassword(r.Password) {
		return echo.ErrUnauthorized
	}
	t, err := a.GenerateJWT()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, NewAuthResponse(t))
}
