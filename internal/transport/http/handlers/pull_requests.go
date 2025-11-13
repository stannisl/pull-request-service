package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mostanin/avito-test/internal/domain"
	"github.com/mostanin/avito-test/internal/service"
)

type PullRequestHandler struct {
	service service.PullRequestService
}

func NewPullRequestHandler(service service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{service: service}
}

func (h *PullRequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req createPRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrorUnknown, "invalid JSON payload")
		return
	}

	pr := domain.PullRequest{
		ID:       req.PullRequestID,
		Name:     req.PullRequestName,
		AuthorID: req.AuthorID,
	}

	created, err := h.service.Create(r.Context(), pr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrorUnknown, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"pr": fromDomainPullRequest(created),
	})
}

func (h *PullRequestHandler) Merge(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req mergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrorUnknown, "invalid JSON payload")
		return
	}

	merged, err := h.service.Merge(r.Context(), domain.PRID(req.PullRequestID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrorUnknown, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"pr": fromDomainPullRequest(merged),
	})
}

func (h *PullRequestHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req reassignPRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrorUnknown, "invalid JSON payload")
		return
	}

	pr, replacedBy, err := h.service.ReassignReviewer(r.Context(), domain.PRID(req.PullRequestID), domain.UserID(req.OldUserID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrorUnknown, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"pr":          fromDomainPullRequest(pr),
		"replaced_by": replacedBy,
	})
}

type createPRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type mergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type reassignPRRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

type pullRequestResponse struct {
	ID                string     `json:"pull_request_id"`
	Name              string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         time.Time  `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

func fromDomainPullRequest(pr domain.PullRequest) pullRequestResponse {
	return pullRequestResponse{
		ID:                pr.ID,
		Name:              pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}
