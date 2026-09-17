package service

import (
	"CommerceCore/internal/users/domain"
	"CommerceCore/internal/users/mocks"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func Test_RegisterSuccess(t *testing.T) {
	mockRepo := new(mocks.MockStorage)
	mockCart := new(mocks.MockCartService)
	model := &domain.User{}
	svc := NewUserServiceImpl(mockRepo, mockCart, "test-secret")
	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(model, nil)

	_, err := svc.Register(context.Background(), model)

	assert.Equal(t, nil, err)
	assert.NoErrorf(t, err, "Register error")
	mockRepo.AssertExpectations(t)
}

func Test_RegisterFail(t *testing.T) {
	mockRepo := new(mocks.MockStorage)
	mockCart := new(mocks.MockCartService)
	model := &domain.User{}
	svc := NewUserServiceImpl(mockRepo, mockCart, "test-secret")
	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(model, domain.FailedCreatedUser)

	_, err := svc.Register(context.Background(), model)

	assert.Equal(t, domain.FailedCreatedUser, err)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func Test_LoginSuccess(t *testing.T) {
	email := "example@mail.ru"
	password := "qwerty1234"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
	}
	mockRepo := new(mocks.MockStorage)
	mockCart := new(mocks.MockCartService)
	svc := NewUserServiceImpl(mockRepo, mockCart, "test-secret")
	mockRepo.On("GetUserByEmail", mock.Anything, email).Return(domain.User{Password: string(hash)}, nil)
	mockCart.On("MergeGuestCart", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err = svc.Login(context.Background(), email, password, "1")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func Test_LoginFail(t *testing.T) {
	email := "example@mail.ru"
	realPassword := "qwerty1234"
	hash, err := bcrypt.GenerateFromPassword([]byte(realPassword), bcrypt.MinCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
	}
	mockRepo := new(mocks.MockStorage)
	mockCart := new(mocks.MockCartService)
	svc := NewUserServiceImpl(mockRepo, mockCart, "test-secret")
	mockRepo.On("GetUserByEmail", mock.Anything, email).Return(domain.User{Password: string(hash)}, nil)

	_, err = svc.Login(context.Background(), email, "wrongpassword", "1")

	assert.Equal(t, domain.InvalidPassword, err)
	mockRepo.AssertExpectations(t)
}

func Test_GetUserByIdSuccess(t *testing.T) {
	mockRepo := new(mocks.MockStorage)
	mockCart := new(mocks.MockCartService)
	svc := NewUserServiceImpl(mockRepo, mockCart, "test-secret")
	mockRepo.On("GetUserByID", mock.Anything, mock.Anything).Return(&domain.User{}, nil)

	_, err := svc.GetById(context.Background(), "1")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func Test_GetByIdFail(t *testing.T) {
	mockRepo := new(mocks.MockStorage)
	mockCart := new(mocks.MockCartService)
	svc := NewUserServiceImpl(mockRepo, mockCart, "test-secret")
	mockRepo.On("GetUserByID", mock.Anything, mock.Anything).Return(&domain.User{}, domain.FailedToGetUser)

	_, err := svc.GetById(context.Background(), "1")

	assert.Equal(t, domain.FailedToGetUser, err)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func Test_UpdateUserSuccess(t *testing.T) {
	mockRepo := new(mocks.MockStorage)
	mockCart := new(mocks.MockCartService)
	svc := NewUserServiceImpl(mockRepo, mockCart, "test-secret")
	mockRepo.On("GetUserByID", mock.Anything, mock.Anything).Return(&domain.User{}, nil)
	mockRepo.On("UpdateUser", mock.Anything, mock.Anything).Return(nil)

	err := svc.UpdateUser(context.Background(), domain.User{})

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

}

func Test_UpdateUserFail(t *testing.T) {
	mockRepo := new(mocks.MockStorage)
	mockCart := new(mocks.MockCartService)
	svc := NewUserServiceImpl(mockRepo, mockCart, "test-secret")
	mockRepo.On("GetUserByID", mock.Anything, mock.Anything).Return((*domain.User)(nil), domain.FailedToGetUser)

	err := svc.UpdateUser(context.Background(), domain.User{Id: "1"})

	assert.Equal(t, domain.FailedToGetUser, err)
	mockRepo.AssertExpectations(t)
}

func Test_UpdateUserRoleFoundUserSuccess(t *testing.T) {
	mockRepo := new(mocks.MockStorage)
	mockCart := new(mocks.MockCartService)
	svc := NewUserServiceImpl(mockRepo, mockCart, "test-secret")
	mockRepo.On("GetUserByID", mock.Anything, mock.Anything).Return(&domain.User{}, nil)
	mockRepo.On("UpdateUserRole", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := svc.UpdateRole(context.Background(), "customer", "1")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

}

func Test_UpdateUserRoleFoundUserFail(t *testing.T) {
	mockRepo := new(mocks.MockStorage)
	mockCart := new(mocks.MockCartService)
	svc := NewUserServiceImpl(mockRepo, mockCart, "test-secret")
	mockRepo.On("GetUserByID", mock.Anything, mock.Anything).Return(&domain.User{}, domain.FailedToGetUser)

	err := svc.UpdateRole(context.Background(), "customer", "1")

	assert.Equal(t, domain.FailedToGetUser, err)
	mockRepo.AssertExpectations(t)

}

func Test_UpdateUserRoleSuccess(t *testing.T) {
	mockRepo := new(mocks.MockStorage)
	mockCart := new(mocks.MockCartService)
	svc := NewUserServiceImpl(mockRepo, mockCart, "test-secret")
	mockRepo.On("GetUserByID", mock.Anything, mock.Anything).Return(&domain.User{}, nil)
	mockRepo.On("UpdateUserRole", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := svc.UpdateRole(context.Background(), "customer", "1")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func Test_UpdateUserRoleInvalidRole(t *testing.T) {
	mockRepo := new(mocks.MockStorage)
	mockCart := new(mocks.MockCartService)
	svc := NewUserServiceImpl(mockRepo, mockCart, "test-secret")
	mockRepo.On("GetUserByID", mock.Anything, mock.Anything).Return(&domain.User{}, nil)

	err := svc.UpdateRole(context.Background(), "user", "1")

	assert.Equal(t, domain.InvalidUserRole, err)
	mockRepo.AssertExpectations(t)
}
