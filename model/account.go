package model

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID       uuid.UUID `json:"uuid"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

func NewAccount(email string, plainPassword string) (*Account, error) {
	if len(plainPassword) == 0 {
		return nil, errors.New("password should not be empty")
	}
	password, err := hashPassword(plainPassword)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:       uuid.New(),
		Email:    email,
		Password: password,
	}, nil
}

type jwtClaims struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	jwt.RegisteredClaims
}

func ValidateAccount(u *Account) bool {
	return true
}

func (a *Account) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(plain))
	return err == nil
}

func (a *Account) GenerateJWT() (string, error) {
	claims := &jwtClaims{
		a.ID,
		a.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return t, nil
}

func hashPassword(plain string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(h), err
}
