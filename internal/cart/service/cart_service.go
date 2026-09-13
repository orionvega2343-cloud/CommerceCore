package service

import (
	"CommerceCore/internal/cart/domain"
	"CommerceCore/internal/cart/domain/errs"
	catalog "CommerceCore/internal/catalog/domain"
	"context"
	"log/slog"
	"slices"
	"time"
)

type CartServiceImpl struct {
	postgres domain.CartRepo
	redis    domain.RedisRepo
	repo     domain.CartItemRepo
	catalog  catalog.ProductService
}

func NewCartService(postgres domain.CartRepo, redis domain.RedisRepo, repo domain.CartItemRepo, catalog catalog.ProductService) *CartServiceImpl {
	return &CartServiceImpl{postgres: postgres, redis: redis, repo: repo, catalog: catalog}
}

func mapGuestItemToCartItem(guestItem domain.CartItem, cartId int) *domain.CartItem {
	return &domain.CartItem{
		CartId:        cartId,
		ProductId:     guestItem.ProductId,
		Quantity:      guestItem.Quantity,
		PriceSnapshot: guestItem.PriceSnapshot,
	}
}

func (s *CartServiceImpl) CreateOrGet(ctx context.Context, userId string) (*domain.Cart, error) {
	cart, err := s.postgres.CreateOrGet(ctx, userId)
	if err != nil {
		slog.Error("Failed to create or get cart", "error", err)
		return nil, err
	}
	return cart, nil
}

func (s *CartServiceImpl) Create(ctx context.Context, item *domain.CartItem) (*domain.CartItem, error) {
	i, err := s.repo.Create(ctx, item)
	if err != nil {
		slog.Error("Failed to create item", "error", err)
		return nil, err
	}
	return i, nil
}

func (s *CartServiceImpl) Update(ctx context.Context, item *domain.CartItem) error {
	err := s.repo.Update(ctx, item)
	if err != nil {
		slog.Error("Failed to update item", "error", err)
		return err
	}
	return nil
}

func (s *CartServiceImpl) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		slog.Error("Failed to delete item", "error", err)
		return err
	}
	return nil
}

func (s *CartServiceImpl) GetGuestCart(ctx context.Context, guestId string) (*domain.Cart, error) {
	cart, err := s.redis.GetGuestCart(ctx, guestId)
	if err != nil {
		slog.Error("failed to get guest cart", "error", err)
		return nil, err
	}
	return cart, nil
}

func (s *CartServiceImpl) AddItemToCart(ctx context.Context, guestId string, productId int, quantity int) (*domain.Cart, error) {
	//Попытка получить корзину
	guestCart, err := s.redis.GetGuestCart(ctx, guestId)
	if err != nil {
		slog.Error("failed to get guest cart", "error", err)
		return nil, err
	}

	//если корзины не существует, создаем новую доменную
	//модель по умолчанию, со статусом active
	if guestCart == nil {
		guestCart = &domain.Cart{Status: "active"}
	}

	//Флаг для нахождения корзины, по умолчанию false
	found := false
	for i := range guestCart.Items {
		if guestCart.Items[i].ProductId == productId {
			guestCart.Items[i].Quantity += quantity
			found = true
			break
		}
	}

	if !found {
		product, err := s.catalog.GetProductById(ctx, productId)
		if err != nil {
			slog.Error("failed to get product", "error", err)
			return nil, err
		}
		guestCart.Items = append(guestCart.Items, domain.CartItem{
			ProductId:     productId,
			Quantity:      quantity,
			PriceSnapshot: product.Price,
		})
	}

	if err := s.redis.SetGuestCart(ctx, guestId, guestCart, 7*24*time.Hour); err != nil {
		slog.Error("failed to set guest cart", "error", err)
		return nil, err
	}

	return guestCart, nil
}

func (s *CartServiceImpl) UpdateItemIntoCart(ctx context.Context, guestId string, productId int, quantity int) (*domain.Cart, error) {
	//Попытка получить гостевую корзину
	guestCart, err := s.redis.GetGuestCart(ctx, guestId)
	if err != nil {
		slog.Error("failed to get guest cart", "error", err)
		return nil, err
	}

	//Если ее нет, вернем ошибку
	if guestCart == nil {
		slog.Error("failed to get guest cart", "error", err)
		return nil, errs.CartNotFound
	}

	//Флаг для проверки: "Нашелся ли продукт?"
	found := false
	for i := range guestCart.Items {
		if guestCart.Items[i].ProductId == productId {
			guestCart.Items[i].Quantity = quantity
			found = true
			break
		}
	}
	//Если не нашлась, выбрасываем Sentinel error
	if !found {
		return nil, errs.ProductNotFound
	}
	//Сохраняем корзину в Redis
	if err := s.redis.SetGuestCart(ctx, guestId, guestCart, 7*24*time.Hour); err != nil {
		slog.Error("failed to set guest cart", "error", err)
		return nil, err
	}
	return guestCart, nil
}

func (s *CartServiceImpl) DeleteItemIntoCart(ctx context.Context, guestId string, productId int) error {
	//Попытка получить корзину
	guestCart, err := s.redis.GetGuestCart(ctx, guestId)
	if err != nil {
		slog.Error("failed to get guest cart", "error", err)
		return err
	}
	//Если корзины нет, вернуть сразу ошибку
	if guestCart == nil {
		return errs.CartNotFound
	}

	//Флаг для проверки: "Нашелся ли продукт?"
	found := false
	for i := range guestCart.Items {
		if guestCart.Items[i].ProductId == productId {
			found = true
			guestCart.Items = slices.Delete(guestCart.Items, i, i+1)
			break
		}
	}
	//Если не нашелся, вернуть ошибку
	if !found {
		slog.Error("failed to delete item", "error", err)
		return errs.ProductNotFound
	}

	if err = s.redis.SetGuestCart(ctx, guestId, guestCart, 7*24*time.Hour); err != nil {
		slog.Error("failed to set guest cart", "error", err)
		return err
	}
	return nil
}

func (s *CartServiceImpl) MergeGuestCart(ctx context.Context, guestId string, userId string) error {
	//Попытка получить гостевую корзину
	guestCart, err := s.redis.GetGuestCart(ctx, guestId)
	if err != nil {
		slog.Error("failed to get guest cart", "error", err)
		return err
	}

	//Промах
	if guestCart == nil {
		return errs.CartNotFound
	}

	//Создаем Postgres корзину
	cart, err := s.postgres.CreateOrGet(ctx, userId)
	if err != nil {
		slog.Error("failed to create guest cart", "error", err)
		return err
	}
	//Сохраняем все в CartItems

	for _, item := range guestCart.Items {
		_, err = s.repo.Create(ctx, mapGuestItemToCartItem(item, cart.Id))
		if err != nil {
			slog.Error("failed to create cart", "error", err)
			return err
		}
	}
	//Удаляем гостевую корзину
	if err = s.redis.DeleteGuestCart(ctx, guestId); err != nil {
		slog.Error("failed to delete guest cart", "error", err)
		return err
	}
	return nil
}
