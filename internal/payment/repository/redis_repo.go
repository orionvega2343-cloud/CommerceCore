package repository

import (
	"CommerceCore/internal/payment/domain"
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ domain.IdempotencyRepo = (*IdempotencyRepoImpl)(nil)

const idempotencyTTL = time.Minute

type IdempotencyRepoImpl struct {
	redis *redis.Client
}

func NewIdempotencyRepo(redis *redis.Client) *IdempotencyRepoImpl {
	return &IdempotencyRepoImpl{redis: redis}
}

// Claim - атомарно занимает ключ (SETNX) Если он уже занят - это НЕ ошибка,
// это сигнал дубля: достаём то, что реально лежит под ключом, и отдаём как
// cachedPaymentId, чтобы вызывающий код не создавал платёж повторно
// ключ уже занят - либо предыдущая попытка успешно завершилась (тут лежит
// настоящий payment.Id, записанный Store), либо она ещё обрабатывается
func (i *IdempotencyRepoImpl) Claim(ctx context.Context, key string) (cachedPaymentId int, claimed bool, err error) {
	ok, err := i.redis.SetNX(ctx, key, "processing", idempotencyTTL).Result()
	if err != nil {
		slog.Error("failed to claim idempotency key", "error", err)
		return 0, false, err
	}
	if ok {
		return 0, true, nil
	}
	val, err := i.redis.Get(ctx, key).Result()
	if err != nil {
		slog.Error("failed to get cached idempotency value", "error", err)
		return 0, false, err
	}
	cachedPaymentId, err = strconv.Atoi(val)
	if err != nil {
		slog.Error("failed to parse cached payment id", "error", err)
		return 0, false, err
	}
	return cachedPaymentId, false, nil
}

// Store - вызывается ПОСЛЕ успешного Create: перезаписывает placeholder
// реальным id платежа, чтобы повторный запрос с тем же ключом получил его
func (i *IdempotencyRepoImpl) Store(ctx context.Context, key string, paymentId int) error {
	if err := i.redis.Set(ctx, key, strconv.Itoa(paymentId), idempotencyTTL).Err(); err != nil {
		slog.Error("failed to store idempotency key", "error", err)
		return err
	}
	return nil
}

// Release - вызывается, если Create после успешного Claim завершился ошибкой:
// освобождает ключ, чтобы легитимный повтор не был заблокирован до истечения TTL
func (i *IdempotencyRepoImpl) Release(ctx context.Context, key string) error {
	if err := i.redis.Del(ctx, key).Err(); err != nil {
		slog.Error("failed to release idempotency key", "error", err)
		return err
	}
	return nil
}
