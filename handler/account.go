package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/prokobit/auth-service/model"
	account "github.com/prokobit/auth-service/store"
	"golang.org/x/crypto/bcrypt"
)

type AccountHandler struct {
	store account.Store
}

func NewAccountHandler() *AccountHandler {
	return &AccountHandler{account.NewMemoryStore()}
}

type NewAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccountResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Token string    `json:"token"`
}

func (h *AccountHandler) SignUp(c echo.Context) error {
	request := new(NewAccountRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	// validate request
	if len(request.Password) == 0 {
		return errors.New("password should not be empty")
	}
	a, err := request.newAccount()
	if err != nil {
		return err
	}
	if err := h.store.Create(a); err != nil {
		return nil
	}
	r, err := accountResponse(a)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, r)
}
func (h *AccountHandler) Login(c echo.Context) error {
	request := new(LoginRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	a, err := h.store.GetByEmail(request.Email)
	if err != nil {
		return echo.ErrUnauthorized
	}
	if !a.CheckPassword(request.Password) {
		return echo.ErrUnauthorized
	}
	r, err := accountResponse(a)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

func (r *NewAccountRequest) newAccount() (*model.Account, error) {
	pass, err := r.hashPassword()
	if err != nil {
		return nil, err
	}
	return &model.Account{
		ID:       uuid.New(),
		Email:    r.Email,
		Password: pass}, nil
}

func (r *NewAccountRequest) hashPassword() (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	return string(h), err
}

func accountResponse(a *model.Account) (*AccountResponse, error) {
	token, err := a.GenerateJWT()
	if err != nil {
		return nil, err
	}
	return &AccountResponse{
		ID:    a.ID,
		Email: a.Email,
		Token: token,
	}, nil
}
