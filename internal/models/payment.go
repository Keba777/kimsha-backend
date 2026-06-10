package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod string

const (
	PaymentCash     PaymentMethod = "cash"
	PaymentCard     PaymentMethod = "card"
	PaymentTelebirr PaymentMethod = "telebirr"
	PaymentCBEPay   PaymentMethod = "cbepay"
)

type Payment struct {
	ID        uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID  uuid.UUID     `gorm:"type:uuid;not null;index" json:"tenant_id"`
	OrderID   uuid.UUID     `gorm:"type:uuid;not null;index" json:"order_id"`
	Order     *Order        `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	CashierID *uuid.UUID    `gorm:"type:uuid" json:"cashier_id"`
	Amount    float64       `gorm:"type:numeric(10,2);not null" json:"amount"`
	Method    PaymentMethod `gorm:"not null" json:"method"`
	Reference string        `json:"reference"`
	Status    string        `gorm:"default:'completed'" json:"status"`
	Notes     string        `json:"notes"`
	CreatedAt time.Time     `json:"created_at"`
}
