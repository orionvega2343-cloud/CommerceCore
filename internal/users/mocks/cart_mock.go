package mocks

import (
	cart "CommerceCore/internal/cart/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

var _ cart.CartService = (*MockCartService)(nil)

// MockCartService - мок cart.CartService, нужен только чтобы UserServiceImpl.Login
// мог вызвать MergeGuestCart в тестах (реальные операции с корзиной здесь не нужны,
// поэтому большинство методов - просто заглушки под интерфейс).
type MockCartService struct {
	mock.Mock
}

func (m *MockCartService) CreateOrGet(ctx context.Context, userId string) (*cart.Cart, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cart.Cart), args.Error(1)
}

func (m *MockCartService) Create(ctx context.Context, item *cart.CartItem) (*cart.CartItem, error) {
	args := m.Called(ctx, item)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cart.CartItem), args.Error(1)
}

func (m *MockCartService) Update(ctx context.Context, item *cart.CartItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockCartService) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCartService) GetGuestCart(ctx context.Context, guestId string) (*cart.Cart, error) {
	args := m.Called(ctx, guestId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cart.Cart), args.Error(1)
}

func (m *MockCartService) AddItemToCart(ctx context.Context, guestId string, productId int, quantity int) (*cart.Cart, error) {
	args := m.Called(ctx, guestId, productId, quantity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cart.Cart), args.Error(1)
}

func (m *MockCartService) UpdateItemIntoCart(ctx context.Context, guestId string, productId int, quantity int) (*cart.Cart, error) {
	args := m.Called(ctx, guestId, productId, quantity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cart.Cart), args.Error(1)
}

func (m *MockCartService) DeleteItemIntoCart(ctx context.Context, guestId string, productId int) error {
	args := m.Called(ctx, guestId, productId)
	return args.Error(0)
}

func (m *MockCartService) MergeGuestCart(ctx context.Context, guestId string, userId string) error {
	args := m.Called(ctx, guestId, userId)
	return args.Error(0)
}
