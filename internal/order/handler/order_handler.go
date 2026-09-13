package handler

import (
	cartErrs "CommerceCore/internal/cart/domain/errs"
	"CommerceCore/internal/order/domain"
	"CommerceCore/internal/order/domain/errs"
	"CommerceCore/internal/order/dto"
	"CommerceCore/pkg/response"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandlerImpl struct {
	svc domain.OrderService
}

func NewOrderHandler(svc domain.OrderService) *OrderHandlerImpl {
	return &OrderHandlerImpl{svc: svc}
}

func toOrderItemResponse(o domain.OrderItem) *dto.OrderItemResponse {
	return &dto.OrderItemResponse{
		Id:           o.Id,
		OrderId:      o.OrderId,
		ProductId:    o.ProductId,
		Quantity:     o.Quantity,
		PricePerUnit: o.PricePerUnit,
	}
}

func toResponseOrder(o domain.Order) *dto.OrderResponse {
	return &dto.OrderResponse{
		Id:          o.Id,
		UserId:      o.UserId,
		CartId:      o.CartId,
		Status:      o.Status,
		TotalAmount: o.TotalAmount,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

// writeError - маппит ошибку сервиса в HTTP-ответ. Единая точка,
// чтобы каждый хендлер не писал свой switch по ошибкам.
func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errs.ErrForbidden):
		c.JSON(http.StatusForbidden, response.Error{Message: err.Error(), Code: "FORBIDDEN"})
	case errors.Is(err, sql.ErrNoRows):
		c.JSON(http.StatusNotFound, response.Error{Message: "order not found", Code: "ORDER_NOT_FOUND"})
	case errors.Is(err, errs.UnknownStatus):
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "INVALID_STATUS_TRANSITION"})
	case errors.Is(err, cartErrs.FailedCheckedOut):
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "CART_NOT_ACTIVE"})
	case errors.Is(err, cartErrs.InvalidCartItemCap):
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "CART_EMPTY"})
	case errors.Is(err, cartErrs.CartNotFound):
		c.JSON(http.StatusNotFound, response.Error{Message: err.Error(), Code: "CART_NOT_FOUND"})
	default:
		c.JSON(http.StatusInternalServerError, response.Error{Message: err.Error(), Code: "INTERNAL_ERROR"})
	}
}

func toOrderListResponse(orders []*domain.Order) *dto.OrderListResponse {
	res := make([]*dto.OrderResponse, 0, len(orders))
	for _, o := range orders {
		res = append(res, toResponseOrder(*o))
	}
	return &dto.OrderListResponse{Orders: res}
}

func (h *OrderHandlerImpl) Checkout(c *gin.Context) {
	userId := c.GetString("user_id")
	ctx := c.Request.Context()

	order, err := h.svc.Checkout(ctx, userId)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponseOrder(*order))
}

func (h *OrderHandlerImpl) GetOrder(c *gin.Context) {
	userId := c.GetString("user_id")
	role := c.GetString("role")

	orderId := c.Param("orderId")
	ctx := c.Request.Context()

	id, err := strconv.Atoi(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_TO_PARSE_ID"})
		return
	}

	order, err := h.svc.GetOrder(ctx, id, userId, role)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponseOrder(*order))
}

func (h *OrderHandlerImpl) ListOrders(c *gin.Context) {
	userId := c.GetString("user_id")
	role := c.GetString("role")

	limit := c.Query("limit")
	offset := c.Query("offset")
	ctx := c.Request.Context()

	parsedLimit, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_TO_PARSE_LIMIT"})
		return
	}

	parsedOffset, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_TO_PARSE_OFFSET"})
		return
	}

	order, err := h.svc.ListOrders(ctx, role, userId, parsedLimit, parsedOffset)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toOrderListResponse(order))
}

func (h *OrderHandlerImpl) TransitionStatus(c *gin.Context) {
	userId := c.GetString("user_id")
	role := c.GetString("role")
	orderId := c.Param("orderId")
	id, err := strconv.Atoi(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_TO_PARSE_ID"})
		return
	}
	ctx := c.Request.Context()

	var req dto.TransitionStatusRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_TO_PARSE_JSON"})
		return
	}

	tStatus, err := h.svc.TransitionStatus(ctx, id, req.Status, userId, role)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponseOrder(*tStatus))
}
