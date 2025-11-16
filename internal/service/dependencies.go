package service

type Dependencies struct {
	TeamService        TeamService
	UserService        UserService
	PullRequestService PullRequestService
	StatsService       StatsService
}
