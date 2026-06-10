package handlers

import (
	"time"

	"kimsha/internal/services"
	"kimsha/pkg/response"
	"kimsha/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PaymentHandler struct{ svc *services.PaymentService }

func NewPaymentHandler(svc *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

func (h *PaymentHandler) Pay(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid order id")
	}
	var in services.PayOrderInput
	if err := c.BodyParser(&in); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	if errs := validator.Validate(in); errs != nil {
		return c.Status(422).JSON(fiber.Map{"errors": errs})
	}
	payment, err := h.svc.PayOrder(orderID, tenantID(c), in)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, payment)
}

func (h *PaymentHandler) List(c *fiber.Ctx) error {
	from := time.Now().Truncate(24 * time.Hour)
	to := from.Add(24 * time.Hour)
	if raw := c.Query("from"); raw != "" {
		if t, err := time.Parse("2006-01-02", raw); err == nil {
			from = t
		}
	}
	if raw := c.Query("to"); raw != "" {
		if t, err := time.Parse("2006-01-02", raw); err == nil {
			to = t.Add(24 * time.Hour)
		}
	}
	payments, err := h.svc.List(tenantID(c), from, to)
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, payments)
}

func (h *PaymentHandler) CashIn(c *fiber.Ctx) error {
	return h.cashTx(c, "in")
}

func (h *PaymentHandler) CashOut(c *fiber.Ctx) error {
	return h.cashTx(c, "out")
}

func (h *PaymentHandler) OpenShift(c *fiber.Ctx) error {
	return h.cashTx(c, "open")
}

func (h *PaymentHandler) CloseShift(c *fiber.Ctx) error {
	return h.cashTx(c, "close")
}

func (h *PaymentHandler) CashSummary(c *fiber.Ctx) error {
	txs, err := h.svc.TodayCashSummary(tenantID(c))
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, txs)
}

func (h *PaymentHandler) cashTx(c *fiber.Ctx, txType string) error {
	var in services.CashTxInput
	if err := c.BodyParser(&in); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	in.Type = txType
	tx, err := h.svc.CreateCashTx(tenantID(c), in)
	if err != nil {
		return response.Internal(c)
	}
	return response.Created(c, tx)
}
