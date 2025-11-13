package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mostanin/avito-test/internal/domain"
	"github.com/mostanin/avito-test/internal/service"
)

type UserHandler struct {
	userService service.UserService
	prService   service.PullRequestService
}

func NewUserHandler(userService service.UserService, prService service.PullRequestService) *UserHandler {
	return &UserHandler{
		userService: userService,
		prService:   prService,
	}
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req setActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrorUnknown, "invalid JSON payload")
		return
	}

	updated, err := h.userService.SetActiveFlag(r.Context(), domain.UserID(req.UserID), req.IsActive)
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrorUnknown, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user": fromDomainUser(updated),
	})
}

func (h *UserHandler) GetReviewAssignments(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, ErrorUnknown, "user_id is required")
		return
	}

	prs, err := h.prService.ListByReviewer(r.Context(), domain.UserID(userID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrorUnknown, err.Error())
		return
	}

	response := reviewAssignmentsResponse{
		UserID: userID,
	}

	for _, pr := range prs {
		response.PullRequests = append(response.PullRequests, pullRequestShort{
			ID:       pr.ID,
			Name:     pr.Name,
			AuthorID: pr.AuthorID,
			Status:   string(pr.Status),
		})
	}

	writeJSON(w, http.StatusOK, response)
}

type setActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type userResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

func fromDomainUser(user domain.User) userResponse {
	return userResponse{
		UserID:   user.ID,
		Username: user.Username,
		TeamName: user.Team,
		IsActive: user.IsActive,
	}
}

type reviewAssignmentsResponse struct {
	UserID       string             `json:"user_id"`
	PullRequests []pullRequestShort `json:"pull_requests"`
}

type pullRequestShort struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
	Status   string `json:"status"`
}
