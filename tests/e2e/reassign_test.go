// e2e/reassign_test.go
package e2e

import (
	"encoding/json"
	"net/http"

	"github.com/stannisl/pull-request-service/internal/transport/dto"
	"github.com/stannisl/pull-request-service/internal/transport/dto/request"
	"github.com/stannisl/pull-request-service/internal/transport/dto/response"
	"github.com/stretchr/testify/assert"
)

func (suite *E2ETestSuite) TestPullRequest_ReassignReviewer() {
	teamReq := request.AddTeamRequest{
		Name: "reassign-team",
		Members: []dto.User{
			{Id: "r-1", Username: "reviewer1", IsActive: true},
			{Id: "r-2", Username: "reviewer2", IsActive: true},
			{Id: "r-3", Username: "reviewer3", IsActive: true},
			{Id: "r-4", Username: "reviewer4", IsActive: true},
			{Id: "r-5", Username: "reviewer5", IsActive: true},
		},
	}

	resp, err := suite.makeRequest("POST", "/team/add", teamReq)
	suite.Require().NoError(err)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	// Создаем PR
	createPRReq := request.CreatePullRequest{
		PullRequestId:   "pr-reassign",
		PullRequestName: "Reassign test",
		AuthorId:        "r-1",
	}

	resp, err = suite.makeRequest("POST", "/pullRequest/create", createPRReq)
	suite.Require().NoError(err)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var prResponse response.PullRequestInfo
	err = json.NewDecoder(resp.Body).Decode(&prResponse)
	suite.Require().NoError(err)

	originalReviewers := prResponse.Pr.AssignedReviewers
	assert.Len(suite.T(), originalReviewers, 2)

	reassignReq := request.ReassignReviewers{
		PullRequestId: "pr-reassign",
		OldReviewerId: originalReviewers[0],
	}

	resp, err = suite.makeRequest("POST", "/pullRequest/reassign", reassignReq)
	if err == nil && resp.StatusCode == http.StatusOK {
		var reassignedPR response.PullRequestInfo
		err = json.NewDecoder(resp.Body).Decode(&reassignedPR)
		suite.Require().NoError(err)

		assert.NotEqual(suite.T(), originalReviewers, reassignedPR.Pr.AssignedReviewers)
		assert.NotContains(suite.T(), reassignedPR.Pr.AssignedReviewers, originalReviewers[0])
	}
}
