package domain

import (
	"CommerceCore/internal/users/dto"
	"context"
)

type UserService interface {
	Register(ctx context.Context, req dto.UserRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req dto.UserRequest) (string, error)
	GetById(ctx context.Context, id string) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id string, req dto.UserRequest) error
	UpdateRole(ctx context.Context, role, id string) error
}
