package e2e

import (
	"encoding/json"
	"net/http"

	"github.com/stannisl/pull-request-service/internal/transport/dto"
	"github.com/stannisl/pull-request-service/internal/transport/dto/request"
	"github.com/stannisl/pull-request-service/internal/transport/dto/response"
	"github.com/stretchr/testify/assert"
)

func (suite *E2ETestSuite) TestTeam_CreateAndGet() {
	createReq := request.AddTeamRequest{
		Name: "backend-team",
		Members: []dto.User{
			{Id: "user-1", Username: "alice", IsActive: true},
			{Id: "user-2", Username: "bob", IsActive: true},
			{Id: "user-3", Username: "charlie", IsActive: true},
		},
	}

	resp, err := suite.makeRequest("POST", "/team/add", createReq)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createResponse response.Team
	err = json.NewDecoder(resp.Body).Decode(&createResponse)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "backend-team", createResponse.Team.TeamName)
	assert.Len(suite.T(), createResponse.Team.Members, 3)

	// Получаем команду
	resp, err = suite.makeRequest("GET", "/team/get?team_name=backend-team", nil)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var getResponse response.Team
	err = json.NewDecoder(resp.Body).Decode(&getResponse)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "backend-team", getResponse.Team.TeamName)
	assert.Len(suite.T(), getResponse.Team.Members, 3)
}

func (suite *E2ETestSuite) TestTeam_NotFound() {
	resp, err := suite.makeRequest("GET", "/team/get?team_name=nonexistent", nil)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *E2ETestSuite) TestTeam_InvalidRequest() {
	invalidReq := map[string]interface{}{
		"team_name": "invalid-team",
	}

	resp, err := suite.makeRequest("POST", "/team/add", invalidReq)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}
