package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stannisl/pull-request-service/internal/domain"
	"github.com/stannisl/pull-request-service/internal/service"
	"github.com/stannisl/pull-request-service/internal/transport/dto"
	"github.com/stannisl/pull-request-service/internal/transport/dto/request"
	"github.com/stannisl/pull-request-service/internal/transport/dto/response"
)

type PullRequestHandler struct {
	service service.PullRequestService
}

func NewPullRequestHandler(service service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{service: service}
}

func (prh *PullRequestHandler) Merge(c *gin.Context) {
	ctx := c.Request.Context()
	var pullRequestMerge request.PullRequestMerge

	if err := c.ShouldBindJSON(&pullRequestMerge); err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Error: dto.ErrBadRequest("Invalid request")})
		return
	}

	pullRequest, err := prh.service.Merge(ctx, &domain.PullRequest{ID: pullRequestMerge.PullRequestId})

	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			c.JSON(http.StatusNotFound, response.Error{Error: dto.ErrNotFound()})
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error{Error: dto.ErrInternalServer(err)})
		return
	}

	var pullRequestResponse response.PullRequestInfo
	c.JSON(http.StatusOK, pullRequestResponse.FromModel(pullRequest))
}

func (prh *PullRequestHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var createPullRequest request.CreatePullRequest

	if err := c.ShouldBindJSON(&createPullRequest); err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Error: dto.ErrBadRequest("Invalid request")})
		return
	}

	pullRequest, err := prh.service.Create(
		ctx,
		createPullRequest.PullRequestId,
		createPullRequest.PullRequestName,
		createPullRequest.AuthorId,
	)

	if err != nil {
		if errors.Is(err, domain.ErrPRExists) {
			c.JSON(http.StatusConflict, response.Error{Error: dto.ErrPRExists()})
			return
		}

		if errors.Is(err, domain.ErrEntityNotFound) {
			c.JSON(http.StatusNotFound, response.Error{Error: dto.ErrNotFound()})
			return
		}

		c.JSON(http.StatusInternalServerError, response.Error{Error: dto.ErrInternalServer(err)})
		return
	}

	var requestInfo response.PullRequestInfo
	c.JSON(http.StatusCreated, requestInfo.FromModel(pullRequest))
}

func (prh *PullRequestHandler) ReassignReviewers(c *gin.Context) {
	ctx := c.Request.Context()
	var reassignReviewers request.ReassignReviewers

	if err := c.ShouldBindJSON(&reassignReviewers); err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Error: dto.ErrBadRequest("Invalid request")})
		return
	}

	pullRequest, _, err := prh.service.ReassignReviewer(
		ctx,
		reassignReviewers.PullRequestId,
		reassignReviewers.OldReviewerId,
	)

	if err != nil {
		if errors.Is(err, domain.ErrPRMerged) {
			c.JSON(http.StatusConflict, response.Error{Error: dto.ErrPRMerged()})
			return
		}

		if errors.Is(err, domain.ErrNotAssigned) {
			c.JSON(http.StatusConflict, response.Error{Error: dto.ErrNotAssigned()})
			return
		}

		if errors.Is(err, domain.ErrNoCandidate) {
			c.JSON(http.StatusConflict, response.Error{Error: dto.ErrNoCandidate()})
			return
		}

		if errors.Is(err, domain.ErrEntityNotFound) {
			c.JSON(http.StatusNotFound, response.Error{Error: dto.ErrNotFound()})
			return
		}

		c.JSON(http.StatusInternalServerError, response.Error{Error: dto.ErrInternalServer(err)})
		return
	}

	var requestInfo response.PullRequestInfo
	c.JSON(http.StatusCreated, requestInfo.FromModel(pullRequest))
}
