package e2e

import (
	"encoding/json"
	"net/http"

	"github.com/stannisl/pull-request-service/internal/transport/dto"
	"github.com/stannisl/pull-request-service/internal/transport/dto/request"
	"github.com/stannisl/pull-request-service/internal/transport/dto/response"
	"github.com/stretchr/testify/assert"
)

func (suite *E2ETestSuite) TestPullRequest_CreateMergeAndReassign() {
	teamReq := request.AddTeamRequest{
		Name: "dev-team",
		Members: []dto.User{
			{Id: "dev-1", Username: "developer1", IsActive: true},
			{Id: "dev-2", Username: "developer2", IsActive: true},
			{Id: "dev-3", Username: "developer3", IsActive: true},
			{Id: "dev-4", Username: "developer4", IsActive: true},
		},
	}

	resp, err := suite.makeRequest("POST", "/team/add", teamReq)
	suite.Require().NoError(err)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	createPRReq := request.CreatePullRequest{
		PullRequestId:   "pr-001",
		PullRequestName: "Add new feature",
		AuthorId:        "dev-1",
	}

	resp, err = suite.makeRequest("POST", "/pullRequest/create", createPRReq)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var prResponse response.PullRequestInfo
	err = json.NewDecoder(resp.Body).Decode(&prResponse)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "pr-001", prResponse.Pr.PullRequestId)
	assert.Equal(suite.T(), "Add new feature", prResponse.Pr.PullRequestName)
	assert.Equal(suite.T(), "dev-1", prResponse.Pr.AuthorId)
	assert.Len(suite.T(), prResponse.Pr.AssignedReviewers, 2)

	mergeReq := request.PullRequestMerge{
		PullRequestId: "pr-001",
	}

	resp, err = suite.makeRequest("POST", "/pullRequest/merge", mergeReq)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var mergedPR response.PullRequestInfo
	err = json.NewDecoder(resp.Body).Decode(&mergedPR)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "MERGED", mergedPR.Pr.Status)
}

func (suite *E2ETestSuite) TestPullRequest_NotFound() {
	mergeReq := request.PullRequestMerge{
		PullRequestId: "nonexistent-pr",
	}

	resp, err := suite.makeRequest("POST", "/pullRequest/merge", mergeReq)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}
