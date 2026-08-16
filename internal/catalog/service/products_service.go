package service

import (
	"CommerceCore/internal/catalog/domain"
	"CommerceCore/internal/catalog/domain/errs"
	"context"
	"log/slog"
	"time"
)

var _ domain.ProductService = (*ProductServiceImpl)(nil)

type ProductServiceImpl struct {
	repo  domain.ProductRepo
	cache domain.ProductCache
}

func NewProductService(repo domain.ProductRepo, cache domain.ProductCache) *ProductServiceImpl {
	return &ProductServiceImpl{repo: repo, cache: cache}
}

func (p *ProductServiceImpl) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	role, ok := ctx.Value("role").(string)
	if !ok {
		return nil, errs.InvalidRole
	}

	if role != "admin" {
		return nil, errs.InvalidRole
	}

	if err := product.Validate(); err != nil {
		return nil, err
	}

	pr, err := p.repo.CreateProduct(ctx, product)
	if err != nil {
		slog.Error("failed to create product", "error", err)
		return nil, err
	}
	return pr, nil
}

// GetAllProducts - метод реализует cache aside паттерн,
// пытается получить данные из кэша, если их нет,
// то создает, если есть отдает пользователю
func (p *ProductServiceImpl) GetAllProducts(ctx context.Context, isActive *bool, limit, offset int) ([]*domain.Product, error) {
	cached, err := p.cache.GetProductList(ctx, isActive, limit, offset)
	if err != nil {
		slog.Error("failed to get product list", "error", err)
	}

	if cached != nil {
		return cached, nil
	}

	products, err := p.repo.GetAllProducts(ctx, isActive, limit, offset)
	if err != nil {
		slog.Error("failed to get product list", "error", err)
		return nil, err
	}
	//Запись в кэш
	if err := p.cache.SetProductList(ctx, isActive, limit, offset, products, 60*time.Second); err != nil {
		slog.Error("failed to cache product list", "error", err)
	}
	return products, nil
}

func (p *ProductServiceImpl) GetProductById(ctx context.Context, id int) (*domain.Product, error) {
	cached, err := p.cache.GetProduct(ctx, id)
	if err != nil {
		slog.Error("cache error", "error", err)
	}
	if cached != nil {
		return cached, nil
	}
	product, err := p.repo.GetProductById(ctx, id)
	if err != nil {
		return nil, err
	}
	//Запись в кэш
	if err := p.cache.SetProduct(ctx, product, 60*time.Second); err != nil {
		slog.Error("failed to set cache", "error", err)
	}
	return product, nil
}

func (p *ProductServiceImpl) UpdateProduct(ctx context.Context, product *domain.Product) error {
	err := p.repo.UpdateProduct(ctx, product)
	if err != nil {
		slog.Error("failed to update product", "error", err)
		return err
	}

	err = p.cache.DeleteProduct(ctx, product.Id)
	if err != nil {
		slog.Error("failed to delete product", "error", err)
	}
	return nil
}

func (p *ProductServiceImpl) DeleteProduct(ctx context.Context, productId int) error {
	err := p.repo.DeleteProduct(ctx, productId)
	if err != nil {
		slog.Error("failed to delete product", "error", err)
		return err
	}
	err = p.cache.DeleteProduct(ctx, productId)
	if err != nil {
		slog.Error("failed to delete product", "error", err)
	}
	return nil
}
