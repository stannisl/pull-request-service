package response

import "github.com/stannisl/avito-test/internal/domain"

type PullRequestInfo struct {
	Pr struct {
		PullRequestId     string          `json:"pull_request_id"`
		PullRequestName   string          `json:"pull_request_name"`
		AuthorId          domain.UserID   `json:"author_id"`
		Status            string          `json:"status"`
		AssignedReviewers []domain.UserID `json:"assigned_reviewers"`
	} `json:"pr"`
}

func (pri *PullRequestInfo) FromModel(model *domain.PullRequest) *PullRequestInfo {
	pri.Pr.AuthorId = model.AuthorID
	pri.Pr.Status = model.Status
	pri.Pr.PullRequestName = model.Name
	pri.Pr.PullRequestId = model.ID
	pri.Pr.AssignedReviewers = model.AssignedReviewers

	return pri
}
