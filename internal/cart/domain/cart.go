package domain

import (
	"CommerceCore/internal/cart/domain/errs"
	"time"
)

type Cart struct {
	Id        int       `db:"id"`
	UserID    string    `db:"user_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Items     []CartItem
}

func (c *Cart) Checkout() error {
	if c.Status != "active" {
		return errs.FailedCheckedOut
	}
	if len(c.Items) == 0 {
		return errs.InvalidCartItemCap
	}
	c.Status = "checked_out"
	return nil
}
