package repository

import (
	"CommerceCore/internal/catalog/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type ProductCacheImpl struct {
	client *redis.Client
}

func NewProductCache(client *redis.Client) *ProductCacheImpl {
	return &ProductCacheImpl{client: client}
}

// buildListKey - формирует ключ для методов GetProductList,
// SetProductList получает 3 значение ключа, во избежание
// конфликта получения данных от разных пользователей
func buildListKey(isActive *bool, limit, offset int) string {
	activeVal := "all"
	if isActive != nil {
		activeVal = strconv.FormatBool(*isActive)
	}
	return fmt.Sprintf("products:list:active=%s:limit=%d:offset=%d", activeVal, limit, offset)
}

// GetProduct - получает данные из кэша по id,
// при промахе возвращает nil без ошибки,
// при отсутствии данных обращаемся к БД
func (c *ProductCacheImpl) GetProduct(ctx context.Context, id int) (*domain.Product, error) {
	key := fmt.Sprintf("product:%d", id)
	val, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // промах, не ошибка
	}
	if err != nil {
		slog.Error("error getting product list", "error", err)
		return nil, err // реальная ошибка Redis
	}
	var product domain.Product
	if err := json.Unmarshal([]byte(val), &product); err != nil {
		return nil, err
	}
	return &product, nil
}

// SetProduct - сохраняет в кэш данные пришедшие из БД в случае промаха внутри метода GetProduct
func (c *ProductCacheImpl) SetProduct(ctx context.Context, product *domain.Product, ttl time.Duration) error {
	key := fmt.Sprintf("product:%d", product.Id)
	data, err := json.Marshal(product)
	if err != nil {
		slog.Error("error marshalling product", "error", err)
		return err
	}
	err = c.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		slog.Error("error while setting product", "error", err)
		return err
	}
	return nil
}

// GetProductList - аналогично GetProduct
func (c *ProductCacheImpl) GetProductList(ctx context.Context, isActive *bool, limit, offset int) ([]*domain.Product, error) {
	key := buildListKey(isActive, limit, offset)
	val, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		slog.Error("error getting product list", "error", err)
		return nil, err
	}

	var products []*domain.Product
	err = json.Unmarshal([]byte(val), &products)
	if err != nil {
		slog.Error("error unmarshalling product list", "error", err)
		return nil, err
	}
	return products, nil
}

// SetProductList - аналогично SetProduct
func (c *ProductCacheImpl) SetProductList(ctx context.Context, isActive *bool, limit, offset int, products []*domain.Product, ttl time.Duration) error {
	key := buildListKey(isActive, limit, offset)
	data, err := json.Marshal(products)
	if err != nil {
		slog.Error("error marshalling product list", "error", err)
		return err
	}
	err = c.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		slog.Error("error while setting product list", "error", err)
		return err
	}
	return nil
}

func (c *ProductCacheImpl) DeleteProduct(ctx context.Context, id int) error {
	key := fmt.Sprintf("product:%d", id)
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		slog.Error("error while deleting product", "error", err)
		return err
	}
	return nil
}
