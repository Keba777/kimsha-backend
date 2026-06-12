package repository

import (
	"time"

	"kimsha/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository struct{ db *gorm.DB }

func NewPaymentRepository(db *gorm.DB) *PaymentRepository { return &PaymentRepository{db: db} }

func (r *PaymentRepository) Create(p *models.Payment) error {
	return r.db.Create(p).Error
}

func (r *PaymentRepository) List(tenantID uuid.UUID, from, to time.Time) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Where("tenant_id = ? AND created_at >= ? AND created_at <= ?", tenantID, from, to).
		Preload("Order").Order("created_at DESC").Find(&payments).Error
	return payments, err
}

func (r *PaymentRepository) Find(id uuid.UUID) (*models.Payment, error) {
	var p models.Payment
	err := r.db.First(&p, "id = ?", id).Error
	return &p, err
}

// CashTransactions

func (r *PaymentRepository) CreateCashTx(tx *models.CashTransaction) error {
	return r.db.Create(tx).Error
}

func (r *PaymentRepository) TodayCashTransactions(tenantID uuid.UUID) ([]models.CashTransaction, error) {
	var txs []models.CashTransaction
	err := r.db.Where("tenant_id = ? AND shift_date = CURRENT_DATE", tenantID).
		Order("created_at ASC").Find(&txs).Error
	return txs, err
}

// TodayBalance returns the running cash balance for today by summing all cash transactions.
// Inflows: open, in, sale. Outflows: close, out, refund.
func (r *PaymentRepository) TodayBalance(tenantID uuid.UUID) (float64, error) {
	type result struct{ Balance float64 }
	var res result
	err := r.db.Model(&models.CashTransaction{}).
		Where("tenant_id = ? AND shift_date = CURRENT_DATE", tenantID).
		Select("COALESCE(SUM(CASE WHEN type IN ('open','in','sale') THEN amount WHEN type IN ('close','out','refund') THEN -amount ELSE 0 END), 0) AS balance").
		Scan(&res).Error
	return res.Balance, err
}
