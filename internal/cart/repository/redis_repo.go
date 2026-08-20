package repository

import (
	"CommerceCore/internal/cart/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepoImpl struct {
	redis *redis.Client
}

func NewRedisRepo(redis *redis.Client) *RedisRepoImpl {
	return &RedisRepoImpl{redis: redis}
}

func (r *RedisRepoImpl) GetGuestCart(ctx context.Context, guestId string) (*domain.Cart, error) {
	key := fmt.Sprintf("guestCart:%s", guestId)
	//Попытка получить гостевую корзину по ключу
	val, err := r.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil //Промах
	}
	if err != nil {
		slog.Error("failed to get guest cart from redis", "error", err)
		return nil, err
	}
	var cart domain.Cart
	if err = json.Unmarshal([]byte(val), &cart); err != nil {
		slog.Error("failed to unmarshal cart from redis", "error", err)
		return nil, err
	}
	return &cart, nil
}

func (r *RedisRepoImpl) SetGuestCart(ctx context.Context, guestId string, cart domain.Cart, ttl time.Duration) error {
	key := fmt.Sprintf("guestCart:%s", guestId)
	data, err := json.Marshal(cart)
	if err != nil {
		slog.Error("failed to marshal cart to redis", "error", err)
		return err
	}
	if err = r.redis.Set(ctx, key, data, ttl).Err(); err != nil {
		slog.Error("failed to set cart to redis", "error", err)
		return err
	}
	return nil
}

func (r *RedisRepoImpl) DelGuestCart(ctx context.Context, guestId string) error {
	key := fmt.Sprintf("guestCart:%s", guestId)
	if err := r.redis.Del(ctx, key).Err(); err != nil {
		slog.Error("failed to delete cart from redis", "error", err)
		return err
	}
	return nil
}
