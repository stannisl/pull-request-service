package repository

type Dependencies struct {
	PullRequestRepository PullRequestRepository
	TeamRepository        TeamRepository
	UserRepository        UserRepository
	StatsRepository       StatsRepository
}
