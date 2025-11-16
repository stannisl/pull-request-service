package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stannisl/pull-request-service/internal/service"
	"github.com/stannisl/pull-request-service/internal/transport/dto"
	"github.com/stannisl/pull-request-service/internal/transport/dto/response"
)

type StatsHandler struct {
	service service.StatsService
}

func NewStatsHandler(statsService service.StatsService) *StatsHandler {
	return &StatsHandler{service: statsService}
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.service.GetStats(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error{Error: dto.ErrInternalServer(err)})
		return
	}

	var statsResponse response.UserAssignments
	c.JSON(http.StatusOK, statsResponse.FromModel(stats))
}
