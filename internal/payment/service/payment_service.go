package service

import (
	"CommerceCore/internal/infra/kafka"
	"CommerceCore/internal/payment/domain"
	"CommerceCore/internal/payment/domain/errs"
	"CommerceCore/internal/payment/dto"
	"CommerceCore/pkg/transaction"
	"context"
	"log/slog"
)

var _ domain.PaymentService = (*PaymentServiceImpl)(nil)

type PaymentServiceImpl struct {
	repo   domain.PaymentRepo
	d      domain.OrderRepo
	tx     transaction.Transactor
	events domain.EventPublisher
}

func NewPaymentService(repo domain.PaymentRepo, d domain.OrderRepo, tx transaction.Transactor, events domain.EventPublisher) *PaymentServiceImpl {
	return &PaymentServiceImpl{repo: repo, d: d, tx: tx, events: events}
}

func toPaymentPayload(p domain.Payment) dto.PaymentPayload {
	return dto.PaymentPayload{
		PaymentId: p.Id,
		OrderId:   p.OrderId,
		Amount:    p.Amount,
		Method:    p.Method,
	}
}

// Create - создание платежа пользователя, метод - получает заказ,
// проверяет его статус с помощью доменного метода TransitionStatus,
// открывает транзакцию для атомарности операций, создает платеж,
// устанавливает ему статус
// idempotencyKey тут не используется - эта реализация не занимается защитой
// от дублей, это дело IdempotencyServiceImpl, который её оборачивает.
func (s *PaymentServiceImpl) Create(ctx context.Context, payment *domain.Payment, idempotencyKey string) (*domain.Payment, error) {
	order, err := s.d.GetOrderById(ctx, payment.OrderId)
	if err != nil {
		slog.Error("failed to get order", "error", err)
		return nil, err
	}

	if _, err = order.TransitionStatus("paid"); err != nil {
		slog.Error("failed to transition status", "error", err)
		return nil, err
	}

	p := &domain.Payment{
		OrderId: payment.OrderId,
		Amount:  payment.Amount,
		Status:  "succeeded",
		Method:  payment.Method,
	}
	var created *domain.Payment
	if err = s.tx.Transaction(ctx, func(ctx context.Context) error {
		payment, err := s.repo.CreatePayment(ctx, p)
		if err != nil {
			slog.Error("failed to create payment", "error", err)
			return err
		}
		err = s.d.MarkOrderPaid(ctx, payment.OrderId)
		if err != nil {
			slog.Error("failed to mark order paid", "error", err)
			return err
		}
		created = payment
		return nil

	}); err != nil {
		slog.Error("failed to open transaction", "error", err)
		return nil, err
	}
	// публикация - после успешного коммита, вне SQL-транзакции (Kafka не её участник).
	// сбой публикации не должен ронять уже успешно созданный платёж - только логируем.
	payload := toPaymentPayload(*created)
	if err := s.events.Publish(ctx, kafka.TopicPaymentSucceeded, payload); err != nil {
		slog.Error("failed to publish payment.succeeded event", "error", err)
	}

	return created, nil
}

func (s *PaymentServiceImpl) GetPaymentById(ctx context.Context, paymentId int, userId, role string) (*domain.Payment, error) {
	payment, err := s.repo.GetPaymentById(ctx, paymentId)
	if err != nil {
		slog.Error("failed to get payment", "error", err)
		return nil, err
	}

	order, err := s.d.GetOrderById(ctx, payment.OrderId)
	if err != nil {
		slog.Error("failed to get order", "error", err)
		return nil, err
	}

	if role != "admin" && order.UserId != userId {
		return nil, errs.InvalidUserOrAdmin
	}

	return payment, nil
}

func (s *PaymentServiceImpl) ListPayments(ctx context.Context, role, userId string, limit, offset int) ([]*domain.Payment, error) {
	if role != "admin" {
		return s.repo.ListByUserId(ctx, userId, limit, offset)
	}
	return s.repo.ListAllPayments(ctx, limit, offset)
}

// TotalAmountPayments - агрегация по тому же принципу, что и ListPayments:
// admin получает сумму по всем платежам, юзер - только по своим.
func (s *PaymentServiceImpl) TotalAmountPayments(ctx context.Context, role, userId string) (int, error) {
	if role != "admin" {
		return s.repo.TotalAmountByUserId(ctx, userId)
	}
	return s.repo.TotalAmountAll(ctx)
}
