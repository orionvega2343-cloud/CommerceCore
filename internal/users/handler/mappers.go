package handler

import (
	"CommerceCore/internal/users/domain"
	"CommerceCore/internal/users/dto"
	"errors"
	"net/http"
)

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
	case errors.Is(err, domain.InvalidUserRole):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
