package errs

import "fmt"

type CustomError struct {
	Code    ErrorCodeEnum
	Message string
}

func (c *CustomError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", c.Code, c.Message)
}

func New(code ErrorCodeEnum, msg string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: msg,
	}
}

type ErrorCodeEnum int

const (
	CodeUnknown ErrorCodeEnum = iota
	CodeMarketNotFound
	CodeInvalidUserID
	CodeInvalidOrderID
)

var (
	ErrUnknown        = New(CodeUnknown, "unknown error")
	ErrMarketNotFound = New(CodeMarketNotFound, "market not found")
	ErrInvalidUserID  = New(CodeInvalidUserID, "invalid user id")
	ErrInvalidOrderID = New(CodeInvalidOrderID, "invalid order id")
)
