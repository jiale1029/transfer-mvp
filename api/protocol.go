package api

import (
	"strings"
	"time"
)

type CreateAccountRequest struct {
	AccountId      string `json:"account_id"`
	InitialBalance string `json:"initial_balance"`
}

type CreateAccountResponse struct {
	AccountId string `json:"account_id"`
	Balance   string `json:"balance"`
}

type GetAccountBalanceResponse struct {
	AccountId string `json:"account_id"`
	Balance   string `json:"balance"`
}

type GetTransactionResponse struct {
	TransactionId string    `json:"transaction_id"`
	From          string    `json:"from_account"`
	To            string    `json:"to_account"`
	Amount        string    `json:"amount"`
	Status        int       `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type SubmitTransactionRequest struct {
	FromAccountId string `json:"source_account_id"`
	ToAccountId   string `json:"destination_account_id"`
	Amount        string `json:"amount"`
}

type SubmitTransactionResponse struct {
	TransactionId string `json:"transaction_id"`
}

// General success and err response to wrap actual pyaload
type SuccessResponse struct {
	StatusCode int32 `json:"status_code"`
	Response   any   `json:"response"`
}

type ErrResponse struct {
	StatusCode string `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func NewErrResponse(err error) ErrResponse {
	comps := strings.SplitAfterN(err.Error(), ":", 2)
	return ErrResponse{
		StatusCode: comps[0][:len(comps[0])-1],
		StatusMsg:  comps[1],
	}
}
