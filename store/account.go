package account

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/prokobit/auth-service/model"
)

type Store interface {
	Create(*model.Account) error
	GetByID(uuid.UUID) (*model.Account, error)
	GetByEmail(string) (*model.Account, error)
}

type MemoryStore struct {
	Accounts []*model.Account
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		Accounts: []*model.Account{},
	}
}

func (s *MemoryStore) Create(a *model.Account) error {
	f, _ := s.GetByEmail(a.Email)
	if f != nil {
		return echo.ErrForbidden
	}
	s.Accounts = append(s.Accounts, a)
	return nil
}

func (s *MemoryStore) GetByID(id uuid.UUID) (*model.Account, error) {
	for _, a := range s.Accounts {
		if id == a.ID {
			return a, nil
		}
	}
	return nil, echo.ErrNotFound
}

func (s *MemoryStore) GetByEmail(email string) (*model.Account, error) {
	for _, a := range s.Accounts {
		if email == a.Email {
			return a, nil
		}
	}
	return nil, echo.ErrNotFound
}
