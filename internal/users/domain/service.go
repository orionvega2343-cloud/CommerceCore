package domain

import "context"

type UserService interface {
	Register(ctx context.Context, user *User) (*User, error)
	Login(ctx context.Context, email, password, guestId string) (string, error)
	GetById(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, user User) error
	UpdateRole(ctx context.Context, role, id string) error
}
