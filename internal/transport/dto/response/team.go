package response

import (
	"github.com/stannisl/pull-request-service/internal/domain"
	"github.com/stannisl/pull-request-service/internal/transport/dto"
)

type Team struct {
	Team struct {
		TeamName string     `json:"team_name"`
		Members  []dto.User `json:"members"`
	} `json:"team"`
}

func (t *Team) FromModel(model *domain.Team) *Team {
	t.Team.TeamName = model.Name
	t.Team.Members = make([]dto.User, len(model.Members))
	for i, user := range model.Members {
		t.Team.Members[i].FromModel(&user)
	}
	return t
}
