package domain

import (
	"context"
	"time"
)

type RedisRepo interface {
	GetGuestCart(ctx context.Context, guestId string) (*Cart, error)
	SetGuestCart(ctx context.Context, guestId string, item *Cart, ttl time.Duration) error
}
