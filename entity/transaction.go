package entity

import "time"

type TransactionStatus struct {
	Status  int
	Message string
}

var (
	TransactionStatusSuccess    = TransactionStatus{Status: 100, Message: "success"}
	TransactionStatusFail       = TransactionStatus{Status: 200, Message: "fail"}
	TransactionStatusInProgress = TransactionStatus{Status: 300, Message: "in progress"}
)

type Transaction struct {
	Id            int64 `gorm:"primaryKey"`
	TransactionId string
	From          string
	To            string
	AmountDollar  int64
	AmountCent    int64
	Status        int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
