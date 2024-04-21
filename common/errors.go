package common

import "fmt"

type CustomErr interface {
	error
	Code() int32
	Message() string
}

type BizError struct {
	Code    int32
	Message string
	Err     error
}

func (b BizError) Error() string {
	if b.Err != nil {
		return fmt.Sprintf("%v:%v:%v", b.Code, b.Message, b.Err)
	}
	return fmt.Sprintf("%v:%v", b.Code, b.Message)
}

func (b BizError) WithMsg(msg string) error {
	return &BizError{Code: b.Code, Message: fmt.Sprintf("%v:%v", b.Error(), msg)}
}

func (b BizError) Wrap(err error) error {
	return BizError{Code: b.Code, Message: b.Message, Err: err}
}

var (
	ErrMySQL = BizError{Code: 1000, Message: "MySQL error"}
	// Request related
	ErrInvalidRequest = BizError{Code: 10001, Message: "Invalid request"}

	// Account related
	ErrAccountNotFound     = BizError{Code: 30001, Message: "Account not found"}
	ErrAccountAlreadyExist = BizError{Code: 30002, Message: "Account already exists"}

	// Transaction related
	ErrTransactionFailed          = BizError{Code: 40000, Message: "Transaction failed"}
	ErrInsufficientBalance        = BizError{Code: 40001, Message: "Insufficient balance to deduct"}
	ErrSourceAccountNotExist      = BizError{Code: 40002, Message: "Source account does not exist"}
	ErrDestinationAccountNotExist = BizError{Code: 40003, Message: "Destination account does not exist"}
	ErrTransactionNotFound        = BizError{Code: 40004, Message: "Transaction does not exist"}
)
