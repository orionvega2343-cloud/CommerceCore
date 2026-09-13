package service

import (
	"CommerceCore/internal/payment/domain"
	"context"
	"log/slog"
)

var _ domain.PaymentService = (*IdempotencyServiceImpl)(nil)

type IdempotencyServiceImpl struct {
	inner       domain.PaymentService
	idempotency domain.IdempotencyRepo
	repo        domain.PaymentRepo
}

func NewIdempotencyService(inner domain.PaymentService, idempotency domain.IdempotencyRepo, repo domain.PaymentRepo) *IdempotencyServiceImpl {
	return &IdempotencyServiceImpl{inner: inner, idempotency: idempotency, repo: repo}
}

// Create - claim ДО вызова внутреннего Create, если ключ уже занят - это дубль,
// отдаём закешированный платёж и НЕ создаём новый, после вызова - Release при
// ошибке (чтобы легитимный повтор не завис до TTL) или Store при успехе
// (чтобы следующий дубль получил готовый результат, а не создал второй платёж)
func (s *IdempotencyServiceImpl) Create(ctx context.Context, payment *domain.Payment, idempotencyKey string) (*domain.Payment, error) {
	key := domain.BuildIdempotencyKey(idempotencyKey)

	cachedPaymentId, claimed, err := s.idempotency.Claim(ctx, key)
	if err != nil {
		slog.Error("failed to claim idempotency key", "error", err)
		return nil, err
	}
	if !claimed {
		return s.repo.GetPaymentById(ctx, cachedPaymentId)
	}

	created, err := s.inner.Create(ctx, payment, idempotencyKey)
	if err != nil {
		if releaseErr := s.idempotency.Release(ctx, key); releaseErr != nil {
			slog.Error("failed to release idempotency key", "error", releaseErr)
		}
		return nil, err
	}

	if storeErr := s.idempotency.Store(ctx, key, created.Id); storeErr != nil {
		// платёж уже успешно создан - не роняем результат из-за сбоя кэша, только логируем
		slog.Error("failed to store idempotency key after successful payment", "error", storeErr)
	}

	return created, nil
}

func (s *IdempotencyServiceImpl) GetPaymentById(ctx context.Context, paymentId int, userId, role string) (*domain.Payment, error) {
	return s.inner.GetPaymentById(ctx, paymentId, userId, role)
}

func (s *IdempotencyServiceImpl) ListPayments(ctx context.Context, role, userId string, limit, offset int) ([]*domain.Payment, error) {
	return s.inner.ListPayments(ctx, role, userId, limit, offset)
}

func (s *IdempotencyServiceImpl) TotalAmountPayments(ctx context.Context, role, userId string) (int, error) {
	return s.inner.TotalAmountPayments(ctx, role, userId)
}
