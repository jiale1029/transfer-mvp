package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jiale1029/transaction/common"
	"github.com/jiale1029/transaction/entity"
)

// HandleSubmitTransaction creates a transaction
func HandleSubmitTransaction(c *gin.Context) {
	var req SubmitTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrResponse(common.ErrInvalidRequest.Wrap(err)))
		return
	}

	var (
		transactionId string
		err           error
	)

	if strings.HasPrefix(req.Amount, "-") {
		c.JSON(http.StatusBadRequest, NewErrResponse(common.ErrInvalidRequest))
		return
	}

	amount := entity.Money(req.Amount)
	if amount.Dollar() == 0 && amount.Cent() == 0 {
		c.JSON(http.StatusBadRequest, NewErrResponse(common.ErrInvalidRequest))
		return
	}
	if transactionId, err = AccountManager.SubmitTransaction(c.Request.Context(), req.FromAccountId, req.ToAccountId, amount.Dollar(), amount.Cent()); err != nil {
		c.JSON(http.StatusBadRequest, NewErrResponse(err))
		return
	}

	resp := &SubmitTransactionResponse{TransactionId: transactionId}
	c.JSON(http.StatusOK, SuccessResponse{Response: resp})
}

// HandleGetTransaction retrieves a transaction
func HandleGetTransaction(c *gin.Context) {
	txId := c.Param("id")
	if txId == "" {
		c.JSON(http.StatusBadRequest, NewErrResponse(common.ErrInvalidRequest.WithMsg("missing id")))
		return
	}
	tx, err := AccountManager.GetTransaction(c.Request.Context(), txId)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrResponse(err))
		return
	}

	resp := &GetTransactionResponse{
		TransactionId: tx.TransactionId,
		From:          tx.From,
		To:            tx.To,
		Amount:        entity.ToMoneyString(tx.AmountDollar, tx.AmountCent),
		CreatedAt:     tx.CreatedAt,
	}
	c.JSON(http.StatusOK, SuccessResponse{Response: resp})
}

// HandleListTransactions retrieves all transactions available
func HandleListTransactions(c *gin.Context) {
	transactions, err := AccountManager.GetAllTransactions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrResponse(err))
		return
	}

	resp := make([]*GetTransactionResponse, 0, len(transactions))
	for _, transaction := range transactions {
		resp = append(resp, &GetTransactionResponse{
			TransactionId: transaction.TransactionId,
			From:          transaction.From,
			To:            transaction.To,
			Amount:        entity.ToMoneyString(transaction.AmountDollar, transaction.AmountCent),
			Status:        transaction.Status,
			CreatedAt:     transaction.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, SuccessResponse{Response: resp})
}
