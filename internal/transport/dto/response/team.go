package response

import (
	"github.com/stannisl/avito-test/internal/domain"
	"github.com/stannisl/avito-test/internal/transport/dto"
)

type Team struct {
	Name    string     `json:"name"`
	Members []dto.User `json:"members"`
}

func (t *Team) FromModel(model *domain.Team) *Team {
	t.Name = model.Name
	t.Members = make([]dto.User, len(model.Members))
	for i, user := range model.Members {
		t.Members[i].FromModel(&user)
	}
	return t
}
