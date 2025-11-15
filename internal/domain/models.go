package domain

import "time"

type (
	TeamName = string
	UserID   = string
	PRID     = string
)

type User struct {
	Id       UserID   `db:"id"`
	Username string   `db:"username"`
	TeamName TeamName `db:"team_name"`
	IsActive bool     `db:"is_active"`
}

type Team struct {
	Name    string `db:"name"`
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
