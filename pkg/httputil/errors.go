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
	ErrClaimsNotFound  = NewError(ErrCodeUnauthorized, "auth context nil")
	ErrArgsRequired    = NewError(ErrCodeArgsRequired, "arguments required")
	ErrArgsInvalid     = NewError(ErrCodeArgsInvalid, CodeMessage[ErrCodeArgsInvalid])
	ErrNotFound        = NewError(ErrCodeNotFound, CodeMessage[ErrCodeNotFound])
	ErrSystemException = NewError(ErrCodeSystemException, CodeMessage[ErrCodeSystemException])
)

func NewError(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d  %s", e.Code, e.Message)
}

func (e *Error) MergeError(err error) *Error {
	if err != nil {
		e.Message += " " + err.Error()
	}
	return e
}
