package domain

import (
	cart "CommerceCore/internal/cart/domain"
	"context"
)

// CartRepo - то, что order-домену нужно от cart для Checkout.
// Интерфейс объявлен здесь (у потребителя), реализуют его методы
// конкретные репозитории cart-домена. Order их только вызывает.
type CartRepo interface {
	// CreateOrGet - получить активную корзину пользователя (или создать пустую).
	CreateOrGet(ctx context.Context, userID string) (*cart.Cart, error)
	// LoadItems - подгрузить позиции корзины (cart.Items после CreateOrGet пустой).
	LoadItems(ctx context.Context, cartID int) ([]cart.CartItem, error)
	// SetStatus - сохранить новый статус корзины (напр. checked_out) в БД.
	SetStatus(ctx context.Context, cartID int, status string) error
}

// StockRepo - то, что order-домену нужно от catalog для Checkout.
type StockRepo interface {
	// DecrementStock - атомарно списать остаток товара.
	// Недостаточный остаток должен приводить к ошибке (откат всей транзакции).
	DecrementStock(ctx context.Context, productID, qty int) error
}

// EventPublisher - то, что order-домену нужно для публикации доменных событий
// (order.created) во внешнюю шину Payload передаётся как any: order не должен
// знать ни про конкретный DTO-пакет, ни про транспорт (Kafka или что угодно ещё) -
// реализацию подставляет DI
type EventPublisher interface {
	Publish(ctx context.Context, eventType string, payload any) error
}
