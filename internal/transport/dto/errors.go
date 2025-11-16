package dto

import (
	"fmt"

	"github.com/stannisl/pull-request-service/internal/domain"
)

const (
	ErrorCodeBadRequest domain.ErrorCode = "BAD_REQUEST"
	ErrorCodeInternal   domain.ErrorCode = "INTERNAL_SERVER_ERROR"
)

type Error struct {
	Code    domain.ErrorCode `json:"code"`
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewError создает новую ошибку
func NewError(code domain.ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func ErrTeamExists() *Error {
	return NewError(domain.ErrorCodeTeamExists, "team_name already exists")
}

func ErrBadRequest(param string) *Error {
	return NewError(ErrorCodeBadRequest, param)
}

func ErrInternalServer(err error) *Error {
	return NewError(ErrorCodeInternal, err.Error())
}

func ErrPRExists() *Error {
	return NewError(domain.ErrorCodePRExists, "PR id already exists")
}

func ErrPRMerged() *Error {
	return NewError(domain.ErrorCodePRMerged, "cannot reassign on merged PR")
}

func ErrNotAssigned() *Error {
	return NewError(domain.ErrorCodeNotAssigned, "reviewer is not assigned to this PR")
}

func ErrNoCandidate() *Error {
	return NewError(domain.ErrorCodeNoCandidate, "no active replacement candidate in team")
}

func ErrNotFound() *Error {
	return NewError(domain.ErrorCodeNotFound, "resource not found")
}
