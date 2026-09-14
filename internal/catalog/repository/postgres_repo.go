package repository

import (
	"CommerceCore/internal/catalog/domain"
	"CommerceCore/internal/catalog/domain/errs"
	"CommerceCore/pkg/querier"
	"CommerceCore/pkg/transaction"
	"context"
	"database/sql"
	"errors"
	"log/slog"
)

var _ domain.ProductRepo = (*ProductRepoImpl)(nil)

type ProductRepoImpl struct {
	q           querier.Querier
	transaction transaction.Transactor
}

func NewProductRepoImpl(q querier.Querier) *ProductRepoImpl {
	return &ProductRepoImpl{q: q}
}

func (r *ProductRepoImpl) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	err := r.q.GetContext(ctx, &product, `INSERT INTO products(name, price, stock_quantity, is_active) VALUES($1, $2, $3, $4) RETURNING id`, product.Name, product.Price, product.StockQuantity, product.IsActive)
	if err != nil {
		slog.Error("error inserting product into database", "error", err)
		return nil, err
	}
	return product, nil
}

func (r *ProductRepoImpl) GetProductById(ctx context.Context, productId int) (*domain.Product, error) {
	var p domain.Product
	err := r.q.GetContext(ctx, &p, `SELECT id, name, price, stock_quantity, is_active FROM products WHERE id = $1`, productId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ProductNotFound
	}
	if err != nil {
		slog.Error("error getting product from database", "error", err)
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepoImpl) GetAllProducts(ctx context.Context, isActive *bool, limit, offset int) ([]*domain.Product, error) {
	var p []*domain.Product

	if isActive == nil {
		err := r.q.SelectContext(ctx, &p, `SELECT id, name, price, stock_quantity, is_active FROM products LIMIT $1 OFFSET $2`, limit, offset)
		if err != nil {
			slog.Error("error getting products from database", "error", err)
			return nil, err
		}
		return p, nil
	}

	err := r.q.SelectContext(ctx, &p, `SELECT id, name, price, stock_quantity, is_active FROM products WHERE is_active = $1 LIMIT $2 OFFSET $3`, *isActive, limit, offset)
	if err != nil {
		slog.Error("error getting products from database", "error", err)
		return nil, err
	}
	return p, nil
}

func (r *ProductRepoImpl) UpdateProduct(ctx context.Context, product *domain.Product) error {
	res, err := r.q.ExecContext(ctx, `UPDATE products SET name = $1, price = $2, stock_quantity = $3, is_active = $4 WHERE id = $5`, product.Name, product.Price, product.StockQuantity, product.IsActive, product.Id)
	if err != nil {
		slog.Error("error updating product from database", "error", err)
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		slog.Error("error getting rows affected", "error", err)
		return err
	}
	if affected == 0 {
		return errs.ProductNotFound
	}
	return nil
}

func (r *ProductRepoImpl) DeleteProduct(ctx context.Context, productId int) error {
	res, err := r.q.ExecContext(ctx, `UPDATE products SET is_active = false WHERE id = $1`, productId)
	if err != nil {
		slog.Error("error deleting product from database", "error", err)
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		slog.Error("error getting rows affected", "error", err)
		return err
	}
	if affected == 0 {
		return errs.ProductNotFound
	}
	return nil
}

func (r *ProductRepoImpl) DecrementStock(ctx context.Context, productID, qty int) error {
	tx, ok := transaction.ExtractTx(ctx)
	q := r.q
	if ok {
		q = tx
	}
	res, err := q.ExecContext(ctx, `UPDATE products SET stock_quantity = stock_quantity - $1 WHERE id = $2 AND stock_quantity >= $1`, qty, productID)
	if err != nil {
		slog.Error("error decrementing stock", "error", err)
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		slog.Error("error getting rows affected", "error", err)
		return err
	}
	if affected == 0 {
		return errs.InsufficientStock
	}
	return nil
}
