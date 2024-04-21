package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jiale1029/transaction/common"
	"github.com/jiale1029/transaction/dal"
	"github.com/jiale1029/transaction/dal/mysql"
	"github.com/jiale1029/transaction/entity"
)

var accountManager = dal.NewAccountDAO(mysql.Database)

// HandleCreateAccount creates an account if not exist
func HandleCreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrResponse(err))
		return
	}

	if req.AccountId == "" {
		c.JSON(http.StatusBadRequest, NewErrResponse(common.ErrInvalidRequest))
		return
	}

	initialBalance := entity.Money(req.InitialBalance)
	if err := accountManager.CreateAccount(c.Request.Context(), req.AccountId, initialBalance.Dollar(), initialBalance.Cent()); err != nil {
		c.JSON(http.StatusBadRequest, NewErrResponse(err))
		return
	}

	resp := &CreateAccountResponse{AccountId: req.AccountId, Balance: req.InitialBalance}
	c.JSON(http.StatusOK, SuccessResponse{Response: resp})
}

// HandleCreateAccount retrieves an account given an account_id
func HandleGetAccount(c *gin.Context) {
	accountId := c.Param("id")
	if accountId == "" {
		c.JSON(http.StatusBadRequest, NewErrResponse(common.ErrInvalidRequest.WithMsg("missing id")))
		return
	}
	account, err := accountManager.GetAccount(c.Request.Context(), accountId)
	if err != nil {
		return
	}

	resp := &GetAccountBalanceResponse{
		AccountId: account.AccountId,
		Balance:   entity.ToMoneyString(account.BalanceDollar, account.BalanceCent),
	}
	c.JSON(http.StatusOK, SuccessResponse{Response: resp})
}

// HandleListAccounts retrieves a list of accounts
func HandleListAccounts(c *gin.Context) {
	accounts, err := accountManager.GetAllAccounts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrResponse(err))
		return
	}
	resp := make([]*GetAccountBalanceResponse, 0, len(accounts))
	for _, acc := range accounts {
		resp = append(resp, &GetAccountBalanceResponse{
			AccountId: acc.AccountId,
			Balance:   entity.ToMoneyString(acc.BalanceDollar, acc.BalanceCent),
		})
	}
	c.JSON(http.StatusOK, SuccessResponse{StatusCode: 0, Response: resp})
}
