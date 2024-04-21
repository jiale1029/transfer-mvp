package common

import "fmt"

type CustomErr interface {
	error
	Code() int32
	Message() string
}

type BizError struct {
	code    int32
	message string
	err     error
}

func (b BizError) Error() string {
	if b.err != nil {
		return fmt.Sprintf("%v:%v:%v", b.code, b.message, b.err)
	}
	return fmt.Sprintf("%v:%v", b.code, b.message)
}

func (b BizError) WithMsg(msg string) error {
	return &BizError{code: b.code, message: fmt.Sprintf("%v:%v", b.Error(), msg)}
}

func (b BizError) Wrap(err error) error {
	return BizError{code: b.code, message: b.message, err: err}
}

var (
	ErrMySQL = BizError{code: 1000, message: "MySQL error"}
	// Request related
	ErrInvalidRequest = BizError{code: 10001, message: "Invalid request"}

	// Account related
	ErrAccountNotFound     = BizError{code: 30001, message: "Account not found"}
	ErrAccountAlreadyExist = BizError{code: 30002, message: "Account already exists"}

	// Transaction related
	ErrTransactionFailed          = BizError{code: 40000, message: "Transaction failed"}
	ErrInsufficientBalance        = BizError{code: 40001, message: "Insufficient balance to deduct"}
	ErrSourceAccountNotExist      = BizError{code: 40002, message: "Source account does not exist"}
	ErrDestinationAccountNotExist = BizError{code: 40003, message: "Destination account does not exist"}
	ErrTransactionNotFound        = BizError{code: 40004, message: "Transaction does not exist"}
)
