package domain

import "time"

type (
	TeamName = string
	UserID   = string
	PRID     = string
)

type User struct {
	Id       UserID
	Username string
	IsActive bool
}

type Team struct {
	Name    string
	Members []User
}

type PullRequestStatus = string

const (
	PullRequestStatusOpen   PullRequestStatus = "OPEN"
	PullRequestStatusMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	ID                PRID
	Name              string
	AuthorID          UserID
	Status            PullRequestStatus
	AssignedReviewers []User
	needMoreReviewers bool
	createdAt         time.Time
	mergedAt          *time.Time
}
