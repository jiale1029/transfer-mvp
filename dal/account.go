package dal

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jiale1029/transaction/common"
	"github.com/jiale1029/transaction/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountDAO struct {
	cli *gorm.DB
}

func NewAccountDAO(db *gorm.DB) *AccountDAO {
	return &AccountDAO{cli: db}
}

// CreateAccount creates an account if the account does not exist
func (m *AccountDAO) CreateAccount(ctx context.Context, accountId string, balanceDollar, balanceCent int64) error {
	err := m.cli.Create(&entity.Account{
		AccountId:     accountId,
		BalanceDollar: balanceDollar,
		BalanceCent:   balanceCent,
	}).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") || strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return common.ErrAccountAlreadyExist
		}
		return common.ErrMySQL.Wrap(err)
	}
	return nil
}

func (m *AccountDAO) GetAccount(ctx context.Context, accountId string) (*entity.Account, error) {
	res := entity.Account{}
	err := m.cli.Where("account_id = ?", accountId).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrAccountNotFound
		}
		return nil, common.ErrMySQL.Wrap(err)
	}
	return &res, nil
}

func (m *AccountDAO) GetAllAccounts(ctx context.Context) ([]*entity.Account, error) {
	var results []*entity.Account
	if err := m.cli.Where("1=1").Find(&results).Error; err != nil {
		return nil, common.ErrMySQL.Wrap(err)
	}
	return results, nil
}

func (m *AccountDAO) SubmitTransaction(ctx context.Context, from, to string, amountDollar, amountCent int64) (string, error) {
	transactionId := uuid.NewString()
	transaction := &entity.Transaction{
		TransactionId: transactionId,
		From:          from,
		To:            to,
		AmountDollar:  amountDollar,
		AmountCent:    amountCent,
		Status:        entity.TransactionStatusInProgress.Status,
	}

	// 1. Insert a transaction record first
	if err := m.cli.Create(&transaction).Error; err != nil {
		return "", common.ErrMySQL.Wrap(err)
	}

	// 2. Start a transaction to initiate transfer
	var fromAccount, toAccount entity.Account
	err := m.cli.Transaction(func(tx *gorm.DB) error {

		innerErr := tx.Transaction(func(tx2 *gorm.DB) error {
			// 2.1 lock from account
			if err := tx2.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("account_id = ?", from).
				First(&fromAccount).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return common.ErrSourceAccountNotExist
				}
				return common.ErrMySQL.Wrap(err)
			}

			// 2.2 lock to account
			if err := tx2.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("account_id = ?", to).
				Find(&toAccount).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return common.ErrDestinationAccountNotExist
				}
				return common.ErrMySQL.Wrap(err)
			}

			// 2.3 check if from account has sufficient balance
			if fromAccount.BalanceDollar < transaction.AmountDollar ||
				(fromAccount.BalanceDollar == transaction.AmountDollar &&
					fromAccount.BalanceCent < transaction.AmountCent) {
				return common.ErrInsufficientBalance
			}

			// 2.4 calculate deduction
			centToDeduct := transaction.AmountCent
			fromAccCent := fromAccount.BalanceCent
			fromAccDollar := fromAccount.BalanceDollar
			for centToDeduct > 0 {
				if fromAccCent >= centToDeduct {
					fromAccCent -= centToDeduct
					centToDeduct = 0
				} else {
					fromAccDollar -= 1
					fromAccCent += 10 * entity.CentUnit // 1 dollar is equivalent to 10 * 0.1 cents
				}
			}
			// deduct dollar
			fromAccDollar -= transaction.AmountDollar

			// 2.5 calculate addition
			toAccCent := toAccount.BalanceCent
			toAccDollar := toAccount.BalanceDollar
			// do addition
			toAccDollar += transaction.AmountDollar + ((toAccCent + transaction.AmountCent) / (10 * entity.CentUnit))
			toAccCent = (toAccCent + transaction.AmountCent) % (10 * entity.CentUnit)

			// 2.6 Update the balance of source account
			if err := tx2.Select("balance_dollar", "balance_cent").Where("account_id = ?", fromAccount.AccountId).Updates(&entity.Account{
				BalanceDollar: fromAccDollar,
				BalanceCent:   fromAccCent,
			}).Error; err != nil {
				return common.ErrMySQL.Wrap(err)
			}

			// 2.7 Update the balance of to account
			if err := tx2.Select("balance_dollar", "balance_cent").Where("account_id = ?", toAccount.AccountId).Updates(&entity.Account{
				BalanceDollar: toAccDollar,
				BalanceCent:   toAccCent,
			}).Error; err != nil {
				return common.ErrMySQL.Wrap(err)
			}
			return nil
		})

		if innerErr != nil {
			transaction.Status = entity.TransactionStatusFail.Status
		} else {
			transaction.Status = entity.TransactionStatusSuccess.Status
		}

		// Update the transaction status
		if err := tx.Where("transaction_id = ?", transaction.Id).Updates(&transaction).Error; err != nil {
			return common.ErrMySQL.Wrap(err)
		}

		if innerErr != nil {
			return innerErr
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return transactionId, nil
}

func (m *AccountDAO) GetTransaction(ctx context.Context, transactionId string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	if err := m.cli.Where("transaction_id = ?", transactionId).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrTransactionNotFound
		}
		return nil, common.ErrMySQL.Wrap(err)
	}
	return &transaction, nil
}

func (m *AccountDAO) GetAllTransactions(ctx context.Context) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	if err := m.cli.Where("1=1").Find(&transactions).Error; err != nil {
		return nil, common.ErrMySQL.Wrap(err)
	}
	return transactions, nil
}
