package domain

import (
	"CommerceCore/internal/order/domain/errs"
	"time"
)

type Order struct {
	Id          int       `db:"id"`
	UserId      string    `db:"user_id"`
	CartId      int       `db:"cart_id"`
	Status      string    `db:"status"`
	TotalAmount int       `db:"total_amount"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// TransitionStatus - бизнес правило смены статуса заказа,
// если статус created - переход возможен только в paid или cancelled,
// если paid - переход в cancelled или shipped, если shipped - переход возможен только в completed,
func (o *Order) TransitionStatus(next string) (string, error) {
	switch o.Status {
	case "created":
		if next == "paid" || next == "cancelled" {
			o.Status = next
			return o.Status, nil
		}
	case "paid":
		if next == "shipped" || next == "cancelled" {
			o.Status = next
			return o.Status, nil
		}
	case "shipped":
		if next == "completed" {
			o.Status = next
			return o.Status, nil
		}
	}
	return "", errs.UnknownStatus
}
