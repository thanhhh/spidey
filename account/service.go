package account

import (
	"context"

	"github.com/satori/go.uuid"
)

type Service interface {
	PostAccount(ctx context.Context, name string) (*Account, error)
	GetAccount(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type accountService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &accountService{repository: r}
}

func (s *accountService) PostAccount(ctx context.Context, name string) (*Account, error) {
	account := &Account{
		ID:   uuid.NewV4().String(),
		Name: name,
	}

	err := s.repository.PutAccount(ctx, *account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *accountService) GetAccount(ctx context.Context, id string) (*Account, error) {
	account, err := s.repository.GetAccountByID(ctx, id)

	return account, err
}

func (s *accountService) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	return s.repository.ListAccounts(ctx, skip, take)
}
