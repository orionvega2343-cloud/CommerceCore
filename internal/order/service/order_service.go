package service

import (
	"CommerceCore/internal/infra/kafka"
	"CommerceCore/internal/order/domain"
	"CommerceCore/internal/order/domain/errs"
	"CommerceCore/internal/order/dto"
	"CommerceCore/pkg/transaction"
	"context"
	"log/slog"
)

type OrderServiceImpl struct {
	repo   domain.OrderRepo
	cart   domain.CartRepo
	stock  domain.StockRepo
	tx     transaction.Transactor
	events domain.EventPublisher
}

func NewOrderService(repo domain.OrderRepo, cart domain.CartRepo, stock domain.StockRepo, tx transaction.Transactor, events domain.EventPublisher) *OrderServiceImpl {
	return &OrderServiceImpl{repo: repo, cart: cart, stock: stock, tx: tx, events: events}
}

// toOrderPayload - собирает узкий payload события order.created из результата Checkout.
// Не отдаём наружу сам domain.Order: у события свой контракт, независимый от того,
// как выглядит внутренняя модель заказа (тот же принцип, что и toResponseOrder в HTTP-слое).
func toOrderPayload(o domain.Order) dto.OrderPayload {
	return dto.OrderPayload{
		OrderId:     o.Id,
		UserId:      o.UserId,
		CartId:      o.CartId,
		TotalAmount: o.TotalAmount,
		Status:      o.Status,
		CreatedAt:   o.CreatedAt,
	}
}

// Checkout - оформление заказа из активной корзины пользователя.
// Читаем корзину и её позиции, проверяем бизнес-правило (active + не пусто),
// затем в одной транзакции: создаём заказ и позиции, списываем остатки товаров
// и переводим корзину в checked_out. Любая ошибка внутри откатывает всё целиком
func (s *OrderServiceImpl) Checkout(ctx context.Context, userId string) (*domain.Order, error) {
	cart, err := s.cart.CreateOrGet(ctx, userId)
	if err != nil {
		slog.Error("failed to get cart", "error", err)
		return nil, err
	}

	items, err := s.cart.LoadItems(ctx, cart.Id)
	if err != nil {
		slog.Error("failed to load items", "error", err)
		return nil, err
	}

	cart.Items = items

	if err = cart.Checkout(); err != nil {
		slog.Error("failed to checkout items", "error", err)
		return nil, err
	}

	var total int
	for _, i := range cart.Items {
		total += i.PriceSnapshot * i.Quantity
	}

	o := domain.Order{
		UserId:      userId,
		CartId:      cart.Id,
		Status:      "created",
		TotalAmount: total,
	}

	var created *domain.Order
	if err = s.tx.Transaction(ctx, func(ctx context.Context) error {
		order, err := s.repo.CreateOrder(ctx, &o)
		if err != nil {
			slog.Error("failed to create order", "error", err)
			return err
		}
		var itms *domain.OrderItem
		for _, item := range cart.Items {
			itms = &domain.OrderItem{
				OrderId:      order.Id,
				ProductId:    item.ProductId,
				Quantity:     item.Quantity,
				PricePerUnit: item.PriceSnapshot,
			}
			oItems, err := s.repo.CreateOrderItem(ctx, itms)
			if err != nil {
				slog.Error("failed to create order item", "error", err)
				return err
			}
			// TODO: DecrementStock(ctx, productID, qty int) error — атомарное списание остатка
			// для order/domain.StockRepo (Checkout). Нужна проверка достаточности остатка
			// (UPDATE ... WHERE stock_quantity >= qty, 0 affected rows → ошибка) и участие
			// в транзакции через transaction.ExtractTx, как в order_repo.go.
			err = s.stock.DecrementStock(ctx, item.ProductId, oItems.Quantity)
			if err != nil {
				slog.Error("failed to decrement stock", "error", err)
				return err
			}
		}
		err = s.cart.SetStatus(ctx, cart.Id, cart.Status)
		if err != nil {
			slog.Error("failed to set status", "error", err)
			return err
		}
		created = order
		return nil
	}); err != nil {
		slog.Error("failed to checkout order", "error", err)
		return nil, err
	}

	// публикация - ПОСЛЕ успешного коммита, вне SQL-транзакции (Kafka, !ее участник)
	// сбой публикации не должен ронять уже успешно оформленный заказ - только логируем.
	payload := toOrderPayload(*created)
	if err := s.events.Publish(ctx, kafka.TopicOrderCreated, payload); err != nil {
		slog.Error("failed to publish order.created event", "error", err)
	}

	return created, nil
}

func (s *OrderServiceImpl) ListOrders(ctx context.Context, role, userId string, limit, offset int) ([]*domain.Order, error) {
	if role == "admin" {
		return s.repo.ListAllOrders(ctx, limit, offset)
	}
	return s.repo.ListOrderByUserId(ctx, userId, limit, offset)
}

func (s *OrderServiceImpl) GetOrder(ctx context.Context, orderId int, userId, role string) (*domain.Order, error) {
	order, err := s.repo.GetOrderById(ctx, orderId)
	if err != nil {
		slog.Error("failed to get order", "error", err)
		return nil, err
	}
	if role != "admin" && order.UserId != userId {
		return nil, errs.ErrForbidden
	}
	return order, nil
}

func (s *OrderServiceImpl) TransitionStatus(ctx context.Context, orderId int, next, userId, role string) (*domain.Order, error) {
	order, err := s.repo.GetOrderById(ctx, orderId)
	if err != nil {
		slog.Error("failed to get order", "error", err)
		return nil, err
	}

	if role != "admin" && (next != "cancelled" || order.UserId != userId) {
		return nil, errs.ErrForbidden
	}

	newStatus, err := order.TransitionStatus(next)
	if err != nil {
		slog.Error("failed to transition status", "error", err)
		return nil, err
	}
	if err = s.repo.UpdateStatus(ctx, orderId, newStatus); err != nil {
		slog.Error("failed to update status", "error", err)
		return nil, err
	}
	return order, nil
}
