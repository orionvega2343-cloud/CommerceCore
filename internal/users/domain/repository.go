package domain

import "context"

type UserRepo interface {
	CreateUser(ctx context.Context, m *User) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	UpdateUser(ctx context.Context, m User) error
	UpdateUserRole(ctx context.Context, role, id string) error
}
