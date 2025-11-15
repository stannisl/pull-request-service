package request

import (
	"github.com/stannisl/avito-test/internal/domain"
)

type CreatePullRequest struct {
	PullRequestId   domain.PRID   `json:"pull_request_id" binding:"required"`
	PullRequestName string        `json:"pull_request_name" binding:"required"`
	AuthorId        domain.UserID `json:"author_id" binding:"required"`
}

func (pr *CreatePullRequest) ToModel() *domain.PullRequest {
	return &domain.PullRequest{
		ID:       pr.PullRequestId,
		Name:     pr.PullRequestName,
		AuthorID: pr.AuthorId,
	}
}

type PullRequestMerge struct {
	PullRequestId domain.PRID `json:"pull_request_id" binding:"required"`
}

type ReassignReviewers struct {
	PullRequestId domain.PRID   `json:"pull_request_id"`
	OldReviewerId domain.UserID `json:"old_reviewer_id"`
}
