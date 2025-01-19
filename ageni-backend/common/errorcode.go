package common

import (
	"fmt"
)

type CustomError struct {
	Code    int32
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

var (
	ErrParameter = &CustomError{
		Code:    20002,
		Message: "parameter error",
	}
	ErrWalletSignature = &CustomError{
		Code:    20802,
		Message: "wallet signature error",
	}
	ErrTokenEmpty = &CustomError{
		Code:    20006,
		Message: "token empty",
	}
	ErrParamCheck = &CustomError{
		Code:    100001,
		Message: "parameter check error",
	}
	ErrGenerateFailed = &CustomError{
		Code:    100002,
		Message: "generate failed",
	}
	ErrGenerating = &CustomError{
		Code:    100003,
		Message: "generating",
	}
)
