package request

import (
	"github.com/stannisl/pull-request-service/internal/domain"
	"github.com/stannisl/pull-request-service/internal/transport/dto"
)

type AddTeamRequest struct {
	Name    string     `json:"team_name" binding:"required"`
	Members []dto.User `json:"members" binding:"required,dive"`
}

func (a *AddTeamRequest) ToModel() domain.Team {
	var team domain.Team

	team.Name = a.Name

	team.Members = make([]domain.User, len(a.Members))
	for i, user := range a.Members {
		m := user.ToModel()
		m.TeamName = a.Name
		team.Members[i] = m
	}

	return team
}
