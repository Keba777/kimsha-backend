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
