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
	ID                PRID              `db:"id"`
	Name              string            `db:"name"`
	AuthorID          UserID            `db:"author_id"`
	Status            PullRequestStatus `db:"status"`
	AssignedReviewers []UserID          `db:"assigned_reviewers"`
	NeedMoreReviewers bool              `db:"need_more_reviewers"`
	CreatedAt         time.Time         `db:"created_at"`
	MergedAt          *time.Time        `db:"merged_at"`
}
