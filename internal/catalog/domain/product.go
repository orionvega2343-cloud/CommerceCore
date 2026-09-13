package domain

import (
	"CommerceCore/internal/catalog/domain/errs"
)

type Product struct {
	Id            int    `db:"id"`
	Name          string `db:"name"`
	Price         int    `db:"price"`
	StockQuantity int    `db:"stock_quantity"`
	IsActive      bool   `db:"is_active"`
}

// Deactivate - реализует мягкое удаление, вместо физического Delete
func (p *Product) Deactivate() {
	p.IsActive = false
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return errs.InvalidProductName
	}
	if p.Price <= 0 {
		return errs.InvalidPrice
	}

	if p.StockQuantity <= 0 {
		return errs.InvalidQuantity
	}
	return nil
}
