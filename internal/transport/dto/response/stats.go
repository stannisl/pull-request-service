package response

import "github.com/stannisl/pull-request-service/internal/domain"

type UserAssignments struct {
	ReviewerID domain.UserID `json:"reviewer_id"`
	Count      int           `json:"assigned_count"`
}

func (u *UserAssignments) FromModel(model []domain.UserAssignments) []UserAssignments {
	response := make([]UserAssignments, 0, len(model))
	for _, assignment := range model {
		response = append(response, UserAssignments{
			ReviewerID: assignment.ReviewerID,
			Count:      assignment.Count,
		})
	}
	return response
}
