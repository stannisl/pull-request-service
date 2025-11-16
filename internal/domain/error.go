package domain

import (
	"errors"
)

// ErrorCode представляет код ошибки согласно OpenAPI спецификации
type ErrorCode string

const (
	ErrorCodeTeamExists  ErrorCode = "TEAM_EXISTS"
	ErrorCodePRExists    ErrorCode = "PR_EXISTS"
	ErrorCodePRMerged    ErrorCode = "PR_MERGED"
	ErrorCodeNotAssigned ErrorCode = "NOT_ASSIGNED"
	ErrorCodeNoCandidate ErrorCode = "NO_CANDIDATE"
	ErrorCodeNotFound    ErrorCode = "NOT_FOUND"
)

var (
	ErrEntityNotFound = errors.New("entity not found in db")

	ErrTeamExists  = errors.New("team already exists")
	ErrPRExists    = errors.New("PR already exists")
	ErrPRMerged    = errors.New("PR merged")
	ErrNotAssigned = errors.New("not assigned")
	ErrNoCandidate = errors.New("no candidate")
	ErrNotFound    = errors.New("not found")
)
