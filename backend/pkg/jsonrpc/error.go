package jsonrpc

import (
	"github.com/omeroid/wdc/backend/pkg/errof"
)

const (
	// ErrorCodeParse is parse error code.
	ErrorCodeParse ErrorCode = -32700
	// ErrorCodeInvalidRequest is invalid request error code.
	ErrorCodeInvalidRequest ErrorCode = -32600
	// ErrorCodeMethodNotFound is method not found error code.
	ErrorCodeMethodNotFound ErrorCode = -32601
	// ErrorCodeInvalidParams is invalid params error code.
	ErrorCodeInvalidParams ErrorCode = -32602
	// ErrorCodeInternal is internal error code.
	ErrorCodeInternal ErrorCode = -32603
	// ErrorCodeServer is server error code.
	ErrorCodeServer ErrorCode = -32000
)

type (
	// A ErrorCode by JSON-RPC 2.0.
	ErrorCode int

	// An Error is a wrapper for a JSON interface value.
	Error struct {
		Code    ErrorCode   `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}
)

// ErrParse returns parse error.
func ErrParse() *Error {
	return &Error{
		Code:    ErrorCodeParse,
		Message: errof.ErrParse.Error(),
	}
}

// ErrInvalidRequest returns invalid request error.
func ErrInvalidRequest() *Error {
	return &Error{
		Code:    ErrorCodeInvalidRequest,
		Message: errof.ErrInvalidRequest.Error(),
	}
}

// ErrMethodNotFound returns method not found error.
func ErrMethodNotFound() *Error {
	return &Error{
		Code:    ErrorCodeMethodNotFound,
		Message: errof.ErrMethodNotFound.Error(),
	}
}

// ErrInvalidParams returns invalid params error.
func ErrInvalidParams() *Error {
	return &Error{
		Code:    ErrorCodeInvalidParams,
		Message: errof.ErrInvalidParams.Error(),
	}
}

// ErrTooLongParameter returns too long parameter error.
func ErrTooLongParameter() *Error {
	return &Error{
		Code:    ErrorCodeInvalidParams,
		Message: errof.ErrTooLongParameter.Error(),
	}
}

// ErrInternal returns internal error.
func ErrInternal() *Error {
	return &Error{
		Code:    ErrorCodeInternal,
		Message: errof.ErrInternal.Error(),
	}
}

// ErrServer returns server error.
func ErrServer(err error) *Error {
	return &Error{
		Code:    ErrorCodeServer,
		Message: err.Error(),
	}
}
