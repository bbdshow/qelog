package httputil

import (
	"fmt"
)

const (
	CodeFailed  = -1
	CodeSuccess = 0

	ErrCodeUnauthorized    = 401
	ErrCodeNotFound        = 404
	ErrCodeSystemException = 500
	ErrCodeArgsInvalid     = 1001
	ErrCodeContextNil      = 1002
	ErrCodeSignException   = 1003
	ErrCodeArgsRequired    = 1004

	ErrCodeOpException = 2000
)

var CodeMessage = map[int]string{
	CodeFailed:  "failed",
	CodeSuccess: "success",

	ErrCodeUnauthorized:    "unauthorized",
	ErrCodeNotFound:        "not found",
	ErrCodeSystemException: "system exception",
	ErrCodeArgsInvalid:     "arguments invalid",
	ErrCodeContextNil:      "context nil",
	ErrCodeSignException:   "sign exception",
	ErrCodeOpException:     "operation exception",
}

var (
	ErrArgsRequired    = Error{Code: ErrCodeArgsRequired, Message: "arguments required"}
	ErrArgsInvalid     = Error{Code: ErrCodeArgsInvalid, Message: CodeMessage[ErrCodeArgsInvalid]}
	ErrClaimsNotFound  = Error{Code: ErrCodeUnauthorized, Message: "auth context nil"}
	ErrNotFound        = Error{Code: ErrCodeNotFound, Message: CodeMessage[ErrCodeNotFound]}
	ErrSignException   = Error{Code: ErrCodeSignException, Message: CodeMessage[ErrCodeSignException]}
	ErrContextNil      = Error{Code: ErrCodeContextNil, Message: CodeMessage[ErrCodeContextNil]}
	ErrSystemException = Error{Code: ErrCodeSystemException, Message: CodeMessage[ErrCodeSystemException]}
)

func New(code int, message string) Error {
	return Error{Code: code, Message: message}
}

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%d  %s", e.Code, e.Message)
}
