package dal

import (
	"context"

	"github.com/jiale1029/transaction/entity"
)

//go:generate mockgen -destination=../mocks/mock_dao.go -source=interface.go
type AccountDAO interface {
	CreateAccount(context.Context, string, int64, int64) error
	GetAccount(context.Context, string) (*entity.Account, error)
	GetAllAccounts(context.Context) ([]*entity.Account, error)
	SubmitTransaction(context.Context, string, string, int64, int64) (string, error)
	GetTransaction(context.Context, string) (*entity.Transaction, error)
	GetAllTransactions(context.Context) ([]*entity.Transaction, error)
}
