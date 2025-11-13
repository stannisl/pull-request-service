package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mostanin/avito-test/internal/domain"
	"github.com/mostanin/avito-test/internal/service"
)

type TeamHandler struct {
	service service.TeamService
}

func NewTeamHandler(service service.TeamService) *TeamHandler {
	return &TeamHandler{service: service}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req teamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrorUnknown, "invalid JSON payload")
		return
	}

	team := req.toDomain()
	created, err := h.service.CreateTeam(r.Context(), team)
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrorUnknown, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"team": fromDomainTeam(created),
	})
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		writeError(w, http.StatusBadRequest, ErrorUnknown, "team_name is required")
		return
	}

	team, err := h.service.GetTeam(r.Context(), teamName)
	if err != nil {
		writeError(w, http.StatusNotFound, ErrorNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, fromDomainTeam(team))
}

type teamRequest struct {
	TeamName string             `json:"team_name"`
	Members  []teamMemberRecord `json:"members"`
}

func (r teamRequest) toDomain() domain.Team {
	members := make([]domain.User, 0, len(r.Members))
	for _, member := range r.Members {
		members = append(members, domain.User{
			ID:       member.UserID,
			Username: member.Username,
			Team:     r.TeamName,
			IsActive: member.IsActive,
		})
	}
	return domain.Team{
		Name:    r.TeamName,
		Members: members,
	}
}

type teamMemberRecord struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type teamResponse struct {
	TeamName string             `json:"team_name"`
	Members  []teamMemberRecord `json:"members"`
}

func fromDomainTeam(team domain.Team) teamResponse {
	members := make([]teamMemberRecord, 0, len(team.Members))
	for _, member := range team.Members {
		members = append(members, teamMemberRecord{
			UserID:   member.ID,
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}
	return teamResponse{
		TeamName: team.Name,
		Members:  members,
	}
}
