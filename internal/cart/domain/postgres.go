package domain

import "context"

type CartRepo interface {
	CreateOrGet(ctx context.Context, userId string) (*Cart, error)
}
