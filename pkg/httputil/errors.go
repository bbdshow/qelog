package httputil

import (
	"fmt"
)

const (
	CodeFailed  = -1
	CodeSuccess = 0

	ErrCodeUnauthorized     = 401
	ErrCodeNotFound         = 404
	ErrCodeSystemException  = 500
	ErrCodeArgsInvalid      = 1001
	ErrCodeContextNil       = 1002
	ErrCodeSignException    = 1003
	ErrCodeSignatureInvalid = 1004
	ErrCodeOpException      = 1005
)

var CodeMessage = map[int]string{
	CodeFailed:  "failed",
	CodeSuccess: "success",

	ErrCodeUnauthorized:     "unauthorized",
	ErrCodeNotFound:         "not found",
	ErrCodeSystemException:  "system exception",
	ErrCodeArgsInvalid:      "arguments invalid",
	ErrCodeContextNil:       "context nil",
	ErrCodeSignException:    "sign exception",
	ErrCodeSignatureInvalid: "signature invalid",
	ErrCodeOpException:      "operation exception",
}

var (
	ErrUnauthorized     = NewError(ErrCodeUnauthorized, CodeMessage[ErrCodeUnauthorized])
	ErrSignatureInvalid = NewError(ErrCodeSignatureInvalid, CodeMessage[ErrCodeSignatureInvalid])
	ErrCtxClaimsNil     = NewError(ErrCodeUnauthorized, "context claims nil")
	ErrArgsInvalid      = NewError(ErrCodeArgsInvalid, CodeMessage[ErrCodeArgsInvalid])
	ErrNotFound         = NewError(ErrCodeNotFound, CodeMessage[ErrCodeNotFound])
	ErrSystemException  = NewError(ErrCodeSystemException, CodeMessage[ErrCodeSystemException])
	ErrOpException      = NewError(ErrCodeOpException, CodeMessage[ErrCodeOpException])
)

func NewError(code int, message string) Error {
	return Error{Code: code, Message: message}
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%d  %s", e.Code, e.Message)
}

func (e Error) MergeError(err error) Error {
	return NewError(e.Code, e.Message+" "+err.Error())
}
func (e Error) MergeString(s string) Error {
	return NewError(e.Code, e.Message+" "+s)
}
