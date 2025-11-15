package domain

import (
	"errors"
	"fmt"
)

// ErrorCode представляет код ошибки согласно OpenAPI спецификации
type ErrorCode string

const (
	// Internal

	ErrorCodeTeamExists  ErrorCode = "TEAM_EXISTS"
	ErrorCodePRExists    ErrorCode = "PR_EXISTS"
	ErrorCodePRMerged    ErrorCode = "PR_MERGED"
	ErrorCodeNotAssigned ErrorCode = "NOT_ASSIGNED"
	ErrorCodeNoCandidate ErrorCode = "NO_CANDIDATE"
	ErrorCodeNotFound    ErrorCode = "NOT_FOUND"

	// Common

	ErrorCodeBadRequest ErrorCode = "BAD_REQUEST"
	ErrorCodeInternal   ErrorCode = "INTERNAL_SERVER_ERROR"
)

var (
	ErrTeamNotFound = errors.New("team not found")
)

type Error struct {
	Code    ErrorCode
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewError создает новую ошибку
func NewError(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func ErrTeamExists(teamName TeamName) *Error {
	return NewError(ErrorCodeTeamExists, fmt.Sprintf("team_name %s already exists", teamName))
}

func ErrBadRequest(param string) *Error {
	return NewError(ErrorCodeBadRequest, param)
}

func ErrInternalServer(err error) *Error {
	return NewError(ErrorCodeInternal, err.Error())
}

func ErrPRExists(prID PRID) *Error {
	return NewError(ErrorCodePRExists, fmt.Sprintf("PR id %s already exists", prID))
}

func ErrPRMerged(prID PRID) *Error {
	return NewError(ErrorCodePRMerged, fmt.Sprintf("cannot reassign on merged PR %s", prID))
}

func ErrNotAssigned(prID PRID, userID UserID) *Error {
	return NewError(ErrorCodeNotAssigned, fmt.Sprintf("reviewer %s is not assigned to PR %s", userID, prID))
}

func ErrNoCandidate(teamName TeamName) *Error {
	return NewError(ErrorCodeNoCandidate, fmt.Sprintf("no active replacement candidate in team %s", teamName))
}

func ErrNotFound(resource string) *Error {
	return NewError(ErrorCodeNotFound, fmt.Sprintf("%s not found", resource))
}
