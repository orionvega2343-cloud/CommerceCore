package mocks

import (
	"CommerceCore/internal/users/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

var _ domain.UserRepo = (*MockStorage)(nil)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) CreateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockStorage) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockStorage) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockStorage) UpdateUser(ctx context.Context, u domain.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockStorage) UpdateUserRole(ctx context.Context, role, id string) error {
	args := m.Called(ctx, role, id)
	return args.Error(0)
}
