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
	TeamName TeamName
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
	AssignedReviewers []UserID // Храним только ID ревьюверов для упрощения
	NeedMoreReviewers bool
	CreatedAt         time.Time
	MergedAt          *time.Time
}
