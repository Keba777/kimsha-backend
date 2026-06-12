package services

import (
	"fmt"
	"time"

	"kimsha/internal/models"
	"kimsha/internal/repository"

	"github.com/google/uuid"
)

type PaymentService struct {
	payRepo   *repository.PaymentRepository
	orderRepo *repository.OrderRepository
	tableRepo *repository.TableRepository
}

func NewPaymentService(
	pr *repository.PaymentRepository,
	or *repository.OrderRepository,
	tr *repository.TableRepository,
) *PaymentService {
	return &PaymentService{payRepo: pr, orderRepo: or, tableRepo: tr}
}

type PayOrderInput struct {
	CashierID *uuid.UUID           `json:"cashier_id"`
	Method    models.PaymentMethod `json:"method" validate:"required"`
	Reference string               `json:"reference"`
}

func (s *PaymentService) PayOrder(orderID, tenantID uuid.UUID, in PayOrderInput) (*models.Payment, error) {
	order, err := s.orderRepo.Find(orderID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("order not found")
	}
	if order.PaymentStatus == models.PaymentStatusPaid {
		return nil, fmt.Errorf("order already paid")
	}

	payment := &models.Payment{
		ID:        uuid.New(),
		TenantID:  tenantID,
		OrderID:   orderID,
		CashierID: in.CashierID,
		Amount:    order.Total,
		Method:    in.Method,
		Reference: in.Reference,
		Status:    "completed",
	}
	if err := s.payRepo.Create(payment); err != nil {
		return nil, err
	}

	now := time.Now()
	order.PaymentStatus = models.PaymentStatusPaid
	order.PaymentMethod = string(in.Method)
	order.Status = models.OrderStatusPaid
	order.PaidAt = &now
	_ = s.orderRepo.Update(order)

	if order.TableID != nil {
		_ = s.tableRepo.UpdateStatus(*order.TableID, tenantID, models.TableFree)
	}

	// For cash payments, record the sale in the cash register so the
	// Cash page balance stays in sync with actual orders.
	if in.Method == models.PaymentCash {
		tableRef := "Takeaway"
		if order.Table != nil {
			tableRef = fmt.Sprintf("Table %d", order.Table.Number)
		}
		note := fmt.Sprintf("Order %s – %s", orderID.String()[:8], tableRef)
		_ = s.recordCashSale(tenantID, in.CashierID, order.Total, note)
	}

	return payment, nil
}

// recordCashSale writes a 'sale' cash transaction and updates the running balance.
func (s *PaymentService) recordCashSale(tenantID uuid.UUID, userID *uuid.UUID, amount float64, note string) error {
	currentBalance, err := s.payRepo.TodayBalance(tenantID)
	if err != nil {
		currentBalance = 0
	}
	tx := &models.CashTransaction{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   userID,
		Type:     "sale",
		Amount:   amount,
		Note:     note,
		Balance:  currentBalance + amount,
	}
	return s.payRepo.CreateCashTx(tx)
}

func (s *PaymentService) List(tenantID uuid.UUID, from, to time.Time) ([]models.Payment, error) {
	return s.payRepo.List(tenantID, from, to)
}

// Cash register

type CashTxInput struct {
	UserID *uuid.UUID `json:"user_id"`
	Type   string     `json:"type" validate:"required,oneof=open close in out"`
	Amount float64    `json:"amount" validate:"required,gt=0"`
	Note   string     `json:"note"`
}

func (s *PaymentService) CreateCashTx(tenantID uuid.UUID, in CashTxInput) (*models.CashTransaction, error) {
	currentBalance, err := s.payRepo.TodayBalance(tenantID)
	if err != nil {
		currentBalance = 0
	}

	var newBalance float64
	switch in.Type {
	case "open", "in":
		newBalance = currentBalance + in.Amount
	case "close", "out":
		newBalance = currentBalance - in.Amount
	default:
		newBalance = currentBalance
	}

	tx := &models.CashTransaction{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   in.UserID,
		Type:     in.Type,
		Amount:   in.Amount,
		Note:     in.Note,
		Balance:  newBalance,
	}
	return tx, s.payRepo.CreateCashTx(tx)
}

func (s *PaymentService) TodayCashSummary(tenantID uuid.UUID) ([]models.CashTransaction, error) {
	return s.payRepo.TodayCashTransactions(tenantID)
}
