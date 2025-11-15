package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stannisl/avito-test/internal/service"
	"github.com/stannisl/avito-test/internal/transport/dto/request"
	"github.com/stannisl/avito-test/internal/transport/dto/response"
)

type PullRequestHandler struct {
	service service.PullRequestService
}

func NewPullRequestHandler(service service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{service: service}
}

func (pr *PullRequestHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var createPullRequest request.CreatePullRequest

	if err := c.ShouldBind(&createPullRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pullRequest, err := pr.service.Create(
		ctx,
		createPullRequest.PullRequestId,
		createPullRequest.PullRequestName,
		createPullRequest.AuthorId,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	var requestInfo response.PullRequestInfo
	c.JSON(http.StatusCreated, requestInfo.FromModel(pullRequest))
}
