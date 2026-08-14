package handler

import (
	"CommerceCore/internal/users/domain"
	"CommerceCore/internal/users/dto"
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

// toDomainUser - собирает доменную модель пользователя из DTO запроса
func toDomainUser(req dto.UserRequest) domain.User {
	return domain.User{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}
}

// toUserResponse - собирает DTO ответа из доменной модели пользователя
func toUserResponse(u *domain.User) dto.UserResponse {
	return dto.UserResponse{
		Id:        u.Id,
		Email:     u.Email,
		Password:  u.Password,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}

func (u *UserHandlerImpl) Register(c *gin.Context) {
	var req dto.UserRequest

	err := c.ShouldBind(&req)
	if err != nil {
		slog.Error("failed to binding type", "error", err)
		c.JSON(http.StatusBadRequest, dto.UserResponse{})
		return
	}
	ctx := c.Request.Context()
	user := toDomainUser(req)
	created, err := u.svc.Register(ctx, &user)
	if err != nil {
		slog.Error("failed to register user", "error", err)
		c.JSON(http.StatusInternalServerError, dto.UserResponse{})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(created))
}

func (u *UserHandlerImpl) Login(c *gin.Context) {
	var req dto.UserRequest

	err := c.ShouldBind(&req)
	if err != nil {
		slog.Error("failed to binding type", "error", err)
		c.JSON(http.StatusBadRequest, dto.UserResponse{})
		return
	}
	ctx := c.Request.Context()
	token, err := u.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		slog.Error("failed to login", "error", err)
		c.JSON(http.StatusInternalServerError, dto.UserResponse{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (u *UserHandlerImpl) GetUserById(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	user, err := u.svc.GetById(ctx, id)
	if err != nil {
		slog.Error("failed to get user by id", "error", err)
		c.JSON(http.StatusInternalServerError, dto.UserResponse{})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user))
}

func (u *UserHandlerImpl) UpdateUser(c *gin.Context) {
	var req dto.UserRequest
	id := c.Param("id")
	err := c.ShouldBind(&req)
	if err != nil {
		slog.Error("failed to binding type", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	ctx := c.Request.Context()
	user := toDomainUser(req)
	user.Id = id
	err = u.svc.UpdateUser(ctx, user)
	if err != nil {
		slog.Error("failed to update user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (u *UserHandlerImpl) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	role := c.Param("role")
	ctx := c.Request.Context()
	err := u.svc.UpdateRole(ctx, role, id)
	if err != nil {
		slog.Error("failed to update role", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}