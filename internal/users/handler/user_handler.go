package handler

import (
	"CommerceCore/internal/users/domain"
	"CommerceCore/internal/users/dto"
	"CommerceCore/pkg/utils"
	"errors"
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

// toUserResponse - собирает DTO ответа из доменной модели пользователя.
// Хэш пароля клиенту никогда не отдаётся.
func toUserResponse(u *domain.User) dto.UserResponse {
	return dto.UserResponse{
		Id:        u.Id,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}

// mapServiceError - переводит доменную ошибку сервиса в HTTP-статус
func mapServiceError(err error) int {
	switch {
	case errors.Is(err, domain.FailedToGetUser):
		return http.StatusNotFound
	case errors.Is(err, domain.InvalidPassword):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
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
		c.JSON(mapServiceError(err), dto.UserResponse{})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(created))
}

func (u *UserHandlerImpl) Login(c *gin.Context) {
	cookie := utils.GuestIdFromCookie(c)
	var req dto.UserRequest
	err := c.ShouldBind(&req)
	if err != nil {
		slog.Error("failed to binding type", "error", err)
		c.JSON(http.StatusBadRequest, dto.UserResponse{})
		return
	}
	ctx := c.Request.Context()
	token, err := u.svc.Login(ctx, req.Email, req.Password, cookie)
	if err != nil {
		slog.Error("failed to login", "error", err)
		c.JSON(mapServiceError(err), dto.UserResponse{})
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
		c.JSON(mapServiceError(err), dto.UserResponse{})
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
		c.JSON(mapServiceError(err), gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (u *UserHandlerImpl) UpdateRole(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateRoleRequest
	err := c.ShouldBind(&req)
	if err != nil {
		slog.Error("failed to binding type", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	ctx := c.Request.Context()
	err = u.svc.UpdateRole(ctx, req.Role, id)
	if err != nil {
		slog.Error("failed to update role", "error", err)
		c.JSON(mapServiceError(err), gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
