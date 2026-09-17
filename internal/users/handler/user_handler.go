package handler

import (
	"CommerceCore/internal/users/domain"
	"CommerceCore/internal/users/dto"
	"CommerceCore/pkg/response"
	"CommerceCore/pkg/utils"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandlerImpl struct {
	svc domain.UserService
}

func NewUserHandlerImpl(svc domain.UserService) *UserHandlerImpl {
	return &UserHandlerImpl{svc: svc}
}

func (u *UserHandlerImpl) Register(c *gin.Context) {
	var req dto.UserRequest

	if err := c.ShouldBind(&req); err != nil {
		slog.Error("failed to binding type", "error", err)
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_TO_BIND"})
		return
	}
	ctx := c.Request.Context()
	user := toDomainUser(req)
	created, err := u.svc.Register(ctx, &user)
	if err != nil {
		slog.Error("failed to register user", "error", err)
		c.JSON(mapServiceError(err), response.Error{Message: err.Error(), Code: "FAILED_TO_REGISTER_USER"})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(created))
}

func (u *UserHandlerImpl) Login(c *gin.Context) {
	cookie := utils.GuestIdFromCookie(c)
	var req dto.UserRequest
	if err := c.ShouldBind(&req); err != nil {
		slog.Error("failed to binding type", "error", err)
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_TO_BIND"})
		return
	}
	ctx := c.Request.Context()
	token, err := u.svc.Login(ctx, req.Email, req.Password, cookie)
	if err != nil {
		slog.Error("failed to login", "error", err)
		c.JSON(mapServiceError(err), response.Error{Message: err.Error(), Code: "FAILED_TO_LOGIN_USER"})
		return
	}
	c.JSON(http.StatusOK, dto.LoginResponse{Token: token})
}

func (u *UserHandlerImpl) GetUserById(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	user, err := u.svc.GetById(ctx, id)
	if err != nil {
		slog.Error("failed to get user by id", "error", err)
		c.JSON(mapServiceError(err), response.Error{Message: err.Error(), Code: "FAILED_TO_GET_USER_BY_ID"})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user))
}

func (u *UserHandlerImpl) UpdateUser(c *gin.Context) {
	var req dto.UserRequest
	id := c.Param("id")
	if err := c.ShouldBind(&req); err != nil {
		slog.Error("failed to binding type", "error", err)
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_TO_BIND"})
		return
	}
	ctx := c.Request.Context()
	user := toDomainUser(req)
	user.Id = id
	if err := u.svc.UpdateUser(ctx, user); err != nil {
		slog.Error("failed to update user", "error", err)
		c.JSON(mapServiceError(err), response.Error{Message: err.Error(), Code: "FAILED_TO_UPDATE_USER"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (u *UserHandlerImpl) UpdateRole(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateRoleRequest
	if err := c.ShouldBind(&req); err != nil {
		slog.Error("failed to binding type", "error", err)
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_TO_BIND"})
		return
	}

	ctx := c.Request.Context()
	if err := u.svc.UpdateRole(ctx, req.Role, id); err != nil {
		slog.Error("failed to update role", "error", err)
		c.JSON(mapServiceError(err), response.Error{Message: err.Error(), Code: "FAILED_TO_UPDATE_USER_ROLE"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
