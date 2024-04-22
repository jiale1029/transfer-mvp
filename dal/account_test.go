package dal

import (
	"context"
	"testing"

	"github.com/jiale1029/transaction/common"
	"github.com/jiale1029/transaction/dal/mysql"
	"github.com/stretchr/testify/assert"
)

func TestAccountDAO_CreateGetAccount(t *testing.T) {
	m := &AccountDAOImpl{cli: mysql.InitMySQL()}

	// create a new account
	err := m.CreateAccount(context.Background(), "123", 123, 500)
	assert.Nil(t, err)

	res, err := m.GetAccount(context.Background(), "123")
	assert.Nil(t, err)
	assert.Equal(t, "123", res.AccountId)
	assert.Equal(t, int64(123), res.BalanceDollar)
	assert.Equal(t, int64(500), res.BalanceCent)

	// unique index
	err = m.CreateAccount(context.Background(), "123", 123, 500)
	assert.Equal(t, common.ErrAccountAlreadyExist.Error(), err.Error())
}

func TestAccountDAO_SubmitTransaction(t *testing.T) {
	m := &AccountDAOImpl{cli: mysql.InitMySQL()}

	// create a new account
	_ = m.CreateAccount(context.Background(), "123", 123, 500)
	_ = m.CreateAccount(context.Background(), "345", 0, 0)

	id, err := m.SubmitTransaction(context.Background(), "123", "345", 35, 10)
	assert.Nil(t, err)

	tx, err := m.GetTransaction(context.Background(), id)
	assert.Nil(t, err)
	assert.Equal(t, id, tx.TransactionId)
	assert.Equal(t, int64(35), tx.AmountDollar)
	assert.Equal(t, int64(10), tx.AmountCent)

	acc1, err := m.GetAccount(context.Background(), "123")
	assert.Nil(t, err)
	acc2, err := m.GetAccount(context.Background(), "345")
	assert.Nil(t, err)

	assert.Equal(t, "123", acc1.AccountId)
	assert.Equal(t, int64(88), acc1.BalanceDollar)
	assert.Equal(t, int64(490), acc1.BalanceCent)
	assert.Equal(t, "345", acc2.AccountId)
	assert.Equal(t, int64(35), acc2.BalanceDollar)
	assert.Equal(t, int64(10), acc2.BalanceCent)
}
