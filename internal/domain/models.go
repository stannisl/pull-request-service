package domain

import "time"

type (
	TeamName = string
	UserID   = string
	PRID     = string
)

type User struct {
	ID       UserID
	Username string
	Team     TeamName
	IsActive bool
}

type Team struct {
	Name    TeamName
	Members []User
}

type PullRequestStatus string

const (
	PullRequestStatusOpen   PullRequestStatus = "OPEN"
	PullRequestStatusMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	ID                PRID
	Name              string
	AuthorID          UserID
	Status            PullRequestStatus
	AssignedReviewers []UserID
	CreatedAt         time.Time
	MergedAt          *time.Time
}
