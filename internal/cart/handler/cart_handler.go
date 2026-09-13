package handler

import (
	"CommerceCore/internal/cart/domain"
	"CommerceCore/internal/cart/domain/errs"
	"CommerceCore/internal/cart/dto"
	"CommerceCore/internal/cart/service"
	"CommerceCore/pkg/response"
	"CommerceCore/pkg/utils"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandlerImpl struct {
	svc service.CartServiceImpl
}

func NewCartHandlerImpl(svc service.CartServiceImpl) *CartHandlerImpl {
	return &CartHandlerImpl{svc: svc}
}

func toDomainCartItem(req *dto.CartItemRequest) *domain.CartItem {
	return &domain.CartItem{
		CartId:        req.CartId,
		ProductId:     req.ProductId,
		Quantity:      req.Quantity,
		PriceSnapshot: req.PriceSnapshot,
	}
}

func toCartItemResponse(c *domain.CartItem) *dto.CartItemResponse {
	return &dto.CartItemResponse{
		Id:            c.Id,
		CartId:        c.CartId,
		ProductId:     c.ProductId,
		Quantity:      c.Quantity,
		PriceSnapshot: c.PriceSnapshot,
	}
}

func toCartResponse(c *domain.Cart) *dto.CartResponse {
	items := make([]dto.CartItemResponse, 0, len(c.Items))
	for _, i := range c.Items {
		items = append(items, dto.CartItemResponse{
			Id:            i.Id,
			ProductId:     i.ProductId,
			Quantity:      i.Quantity,
			PriceSnapshot: i.PriceSnapshot,
		})
	}
	return &dto.CartResponse{
		Id:        c.Id,
		UserID:    c.UserID,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Items:     items,
	}

}

func (h *CartHandlerImpl) CreateOrGet(c *gin.Context) {
	ctx := c.Request.Context()
	userId, ok := utils.UserIdFromContext(c)
	if !ok {
		slog.Error("failed to get userId", "error", errs.FailedGetUserIdFromContext)
		c.JSON(http.StatusBadRequest, response.Error{Message: "failed to get user id from context", Code: "FAILED_GET_TO_CONTEXT"})
		return
	}
	cart, err := h.svc.CreateOrGet(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error{Message: err.Error(), Code: "FAILED_TO_GET_OR_CREATE_CART"})
		return
	}
	c.JSON(http.StatusOK, toCartResponse(cart))
}

func (h *CartHandlerImpl) Create(c *gin.Context) {
	//Получаем userId из контекста
	userId, ok := utils.UserIdFromContext(c)
	if !ok {
		slog.Error("failed to get userId", "error", errs.FailedGetUserIdFromContext)
		c.JSON(http.StatusInternalServerError, response.Error{Message: "failed to get user id from context", Code: "FAILED_GET_TO_CONTEXT"})
		return
	}

	//Парсим JSON в DTO
	var req dto.CartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_PARSE_CART_ITEM"})
		return
	}

	//Получаем пришедшие данные из контекста
	ctx := c.Request.Context()

	//Создаем корзину
	cart, err := h.svc.CreateOrGet(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error{Message: err.Error(), Code: "FAILED_TO_GET_OR_CREATE_CART"})
		return
	}

	//Вручную меняем, CartId пришедший от клиента,
	//для избежания подставления чужого CartId
	item := toDomainCartItem(&req)
	item.CartId = cart.Id

	//Создаем CartItem
	created, err := h.svc.Create(ctx, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error{Message: err.Error(), Code: "FAILED_TO_CREATE_CART_ITEM"})
		return
	}

	c.JSON(http.StatusCreated, toCartItemResponse(created))
}

func (h *CartHandlerImpl) Update(c *gin.Context) {
	id := c.Param("id")
	parseId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_PARSE_ID"})
		return
	}
	var req dto.CartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error{Message: err.Error(), Code: "FAILED_PARSE_CART_ITEM"})
		return
	}
	item := toDomainCartItem(&req)
	item.Id = parseId

	ctx := c.Request.Context()
	if err = h.svc.Update(ctx, item); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error{Message: err.Error(), Code: "FAILED_UPDATE_ITEM"})
		return
	}
	c.JSON(http.StatusOK, toCartItemResponse(item))
}

func (h *CartHandlerImpl) Delete(c *gin.Context) {
	id := c.Param("id")
	parseId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_PARSE_ID"})
		return
	}
	ctx := c.Request.Context()
	if err = h.svc.Delete(ctx, parseId); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error{Message: err.Error(), Code: "FAILED_DELETE_ITEM"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *CartHandlerImpl) GetGuestCart(c *gin.Context) {
	cookie := utils.GuestIdFromCookie(c)
	ctx := c.Request.Context()
	guestCart, err := h.svc.GetGuestCart(ctx, cookie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error{Message: err.Error(), Code: "FAILED_GET_CART_ITEM"})
		return
	}
	c.JSON(http.StatusOK, toCartResponse(guestCart))
}

func (h *CartHandlerImpl) AddItemToCart(c *gin.Context) {
	cookie := utils.GuestIdFromCookie(c)
	ctx := c.Request.Context()

	var req dto.CartItemRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error{Message: err.Error(), Code: "FAILED_PARSE_CART_ITEM"})
		return
	}

	cart, err := h.svc.AddItemToCart(ctx, cookie, req.ProductId, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error{Message: err.Error(), Code: "FAILED_ADD_CART"})
		return
	}
	c.JSON(http.StatusOK, toCartResponse(cart))
}

func (h *CartHandlerImpl) UpdateItemIntoCart(c *gin.Context) {
	cookie := utils.GuestIdFromCookie(c)
	ctx := c.Request.Context()

	var req dto.CartItemRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error{Message: err.Error(), Code: "FAILED_PARSE_CART_ITEM"})
		return
	}

	cart, err := h.svc.UpdateItemIntoCart(ctx, cookie, req.ProductId, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error{Message: err.Error(), Code: "FAILED_UPDATE_ITEM"})
		return
	}
	c.JSON(http.StatusOK, toCartResponse(cart))
}

func (h *CartHandlerImpl) DeleteItemIntoCart(c *gin.Context) {
	cookie := utils.GuestIdFromCookie(c)
	ctx := c.Request.Context()
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_PARSE_ID"})
		return
	}
	err = h.svc.DeleteItemIntoCart(ctx, cookie, productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error{Message: err.Error(), Code: "FAILED_DELETE_ITEM"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}
