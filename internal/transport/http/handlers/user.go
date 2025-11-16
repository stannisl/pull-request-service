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

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (uh *UserHandler) ListPRAsReviewer(c *gin.Context) {
	ctx := c.Request.Context()
	userId := c.Query("user_id")

	if userId == "" {
		c.JSON(http.StatusBadRequest, response.Error{Error: dto.ErrBadRequest("Invalid user_id")})
		return
	}

	prs, err := uh.userService.GetReview(ctx, userId)

	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			c.JSON(http.StatusNotFound, response.Error{Error: dto.ErrNotFound()})
			return
		}

		c.JSON(http.StatusInternalServerError, response.Error{Error: dto.ErrInternalServer(err)})
		return
	}

	var responseUserReviews response.UserReviews
	c.JSON(http.StatusOK, responseUserReviews.MapFrom(userId, prs))
}

func (uh *UserHandler) SetIsActive(c *gin.Context) {
	ctx := c.Request.Context()
	var user request.UserIsActive

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Error: dto.ErrBadRequest("Invalid request")})
		return
	}

	updatedUser, err := uh.userService.SetIsActive(ctx, user.UserId, user.IsActive)

	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			c.JSON(http.StatusNotFound, response.Error{Error: dto.ErrNotFound()})
			return
		}

		c.JSON(http.StatusInternalServerError, response.Error{Error: dto.ErrInternalServer(err)})
		return
	}

	var responseUserInfo response.UserInfo
	c.JSON(http.StatusOK, responseUserInfo.FromModel(updatedUser))
}
