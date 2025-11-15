package response

import "github.com/stannisl/avito-test/internal/domain"

type UserInfo struct {
	User struct {
		UserId   domain.UserID `json:"user_id"`
		Username string        `json:"username"`
		TeamName string        `json:"team_name"`
		IsActive bool          `json:"is_active"`
	} `json:"user"`
}

func (ui *UserInfo) FromModel(user *domain.User) *UserInfo {
	ui.User.UserId = user.Id
	ui.User.Username = user.Username
	ui.User.TeamName = user.TeamName
	ui.User.IsActive = user.IsActive

	return ui
}

type UserReviews struct {
	UserId       domain.UserID        `json:"user_id"`
	PullRequests []PullRequestForUser `json:"pull_requests"`
}

type PullRequestForUser struct {
	PullRequestId   domain.PRID              `json:"pull_request_id"`
	PullRequestName string                   `json:"pull_request_name"`
	AuthorId        domain.UserID            `json:"author_id"`
	Status          domain.PullRequestStatus `json:"status"`
}

func (ur *UserReviews) MapFrom(id domain.UserID, prs []domain.PullRequest) *UserReviews {
	ur.UserId = id

	ur.PullRequests = make([]PullRequestForUser, 0, len(prs))

	for _, pr := range prs {
		ur.PullRequests = append(ur.PullRequests, PullRequestForUser{
			PullRequestId:   pr.ID,
			PullRequestName: pr.Name,
			AuthorId:        pr.AuthorID,
			Status:          pr.Status,
		})
	}

	return ur
}
