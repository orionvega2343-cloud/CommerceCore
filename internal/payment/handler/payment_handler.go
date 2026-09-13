package handler

import (
	orderErrs "CommerceCore/internal/order/domain/errs"
	"CommerceCore/internal/payment/domain"
	"CommerceCore/internal/payment/domain/errs"
	"CommerceCore/internal/payment/dto"
	"CommerceCore/pkg/response"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var _ domain.PaymentHandler = (*PaymentHandlerImpl)(nil)

const (
	defaultLimit  = 20
	defaultOffset = 0
)

type PaymentHandlerImpl struct {
	svc domain.PaymentService
}

func NewPaymentHandler(svc domain.PaymentService) *PaymentHandlerImpl {
	return &PaymentHandlerImpl{svc: svc}
}

// toDomainPayment - req несёт только Amount/Method (то, что реально решает клиент);
// orderId приходит отдельным параметром, потому что берётся из URL, не из тела запроса.
func toDomainPayment(req dto.PaymentRequest, orderId int) *domain.Payment {
	return &domain.Payment{
		OrderId: orderId,
		Amount:  req.Amount,
		Method:  req.Method,
	}
}

func toPaymentResponse(p domain.Payment) *dto.PaymentResponse {
	return &dto.PaymentResponse{
		Id:        p.Id,
		OrderId:   p.OrderId,
		Amount:    p.Amount,
		Status:    p.Status,
		Method:    p.Method,
		CreatedAt: p.CreatedAt,
	}
}

func toPaymentListResponse(payments []*domain.Payment) *dto.PaymentListResponse {
	res := make([]*dto.PaymentResponse, 0, len(payments))
	for _, p := range payments {
		res = append(res, toPaymentResponse(*p))
	}
	return &dto.PaymentListResponse{Payments: res}
}

// writeError - маппит ошибку сервиса в HTTP-ответ, единая точка вместо
// одинакового 400 на любую причину в каждом хендлере (зеркалит order).
func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errs.InvalidUserOrAdmin):
		c.JSON(http.StatusForbidden, response.Error{Message: err.Error(), Code: "FORBIDDEN"})
	case errors.Is(err, orderErrs.OrderNotPayable):
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "ORDER_NOT_PAYABLE"})
	case errors.Is(err, orderErrs.UnknownStatus):
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "INVALID_STATUS_TRANSITION"})
	case errors.Is(err, sql.ErrNoRows):
		c.JSON(http.StatusNotFound, response.Error{Message: "payment not found", Code: "PAYMENT_NOT_FOUND"})
	default:
		c.JSON(http.StatusInternalServerError, response.Error{Message: err.Error(), Code: "INTERNAL_ERROR"})
	}
}

func (p *PaymentHandlerImpl) Create(c *gin.Context) {
	var req dto.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_BIND"})
		return
	}
	orderId, err := strconv.Atoi(c.Param("orderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_PARSE_ORDER_ID"})
		return
	}
	pay := toDomainPayment(req, orderId)
	ctx := c.Request.Context()
	idempotencyKey := c.GetHeader("Idempotency-Key")

	payment, err := p.svc.Create(ctx, pay, idempotencyKey)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toPaymentResponse(*payment))
}

func (p *PaymentHandlerImpl) GetPaymentById(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_PARSE_ID"})
		return
	}
	userId := c.GetString("user_id")
	role := c.GetString("role")

	payment, err := p.svc.GetPaymentById(ctx, id, userId, role)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toPaymentResponse(*payment))
}

func (p *PaymentHandlerImpl) ListPayments(c *gin.Context) {
	ctx := c.Request.Context()
	role := c.GetString("role")
	userId := c.GetString("user_id")

	limit := defaultLimit
	if l := c.Query("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_PARSE_LIMIT"})
			return
		}
		limit = parsed
	}

	offset := defaultOffset
	if o := c.Query("offset"); o != "" {
		parsed, err := strconv.Atoi(o)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_PARSE_OFFSET"})
			return
		}
		offset = parsed
	}

	payments, err := p.svc.ListPayments(ctx, role, userId, limit, offset)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toPaymentListResponse(payments))
}

func (p *PaymentHandlerImpl) TotalAmountPayments(c *gin.Context) {
	ctx := c.Request.Context()
	role := c.GetString("role")
	userId := c.GetString("user_id")

	total, err := p.svc.TotalAmountPayments(ctx, role, userId)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_amount": total})
}
