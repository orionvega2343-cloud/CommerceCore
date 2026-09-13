package domain

import (
	"context"
)

type CartService interface {
	CreateOrGet(ctx context.Context, userId string) (*Cart, error)
	Create(ctx context.Context, item *CartItem) (*CartItem, error)
	Update(ctx context.Context, item *CartItem) error
	Delete(ctx context.Context, id int) error

	GetGuestCart(ctx context.Context, guestId string) (*Cart, error)
	AddItemToCart(ctx context.Context, guestId string, productId int, quantity int) (*Cart, error)
	UpdateItemIntoCart(ctx context.Context, guestId string, productId int, quantity int) (*Cart, error)
	DeleteItemIntoCart(ctx context.Context, guestId string, productId int) error
	MergeGuestCart(ctx context.Context, guestId string, userId string) error
}
