package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stannisl/avito-test/internal/domain"
	"github.com/stannisl/avito-test/internal/service"
	"github.com/stannisl/avito-test/internal/transport/dto/request"
	"github.com/stannisl/avito-test/internal/transport/dto/response"
)

type TeamHandler struct {
	teamService service.TeamService
}

func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

func (t *TeamHandler) AddTeam(c *gin.Context) {
	ctx := c.Request.Context()
	var addTeamRequest request.AddTeamRequest

	if err := c.BindJSON(&addTeamRequest); err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Error: domain.ErrBadRequest("Invalid request")})
		return
	}

	team, err := t.teamService.CreateTeam(ctx, addTeamRequest.ToModel())

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error{Error: domain.ErrInternalServer(err)})
		return
	}

	var teamResponse response.Team
	c.JSON(http.StatusCreated, teamResponse.FromModel(team))
}

func (t *TeamHandler) GetTeam(c *gin.Context) {
	ctx := c.Request.Context()

	teamName := c.Query("team_name")

	if teamName == "" {
		c.JSON(http.StatusBadRequest, response.Error{Error: domain.ErrBadRequest("Invalid request")})
		return
	}

	team, err := t.teamService.GetTeam(ctx, teamName)

	if errors.Is(err, domain.ErrTeamNotFound) {
		c.JSON(http.StatusNotFound, response.Error{Error: domain.ErrNotFound(teamName)})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error{Error: domain.ErrInternalServer(err)})
		return
	}

	var teamResponse response.Team
	c.JSON(http.StatusOK, teamResponse.FromModel(team))
}
