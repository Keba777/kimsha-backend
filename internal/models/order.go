package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderType string
type OrderStatus string
type PaymentStatus string

const (
	OrderTypeDineIn   OrderType = "dine_in"
	OrderTypeTakeaway OrderType = "takeaway"
	OrderTypeDelivery OrderType = "delivery"

	OrderStatusOpen      OrderStatus = "open"
	OrderStatusInKitchen OrderStatus = "in_kitchen"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusServed    OrderStatus = "served"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"

	PaymentStatusUnpaid   PaymentStatus = "unpaid"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type Order struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	TableID        *uuid.UUID     `gorm:"type:uuid;index" json:"table_id"`
	Table          *Table         `gorm:"foreignKey:TableID" json:"table,omitempty"`
	WaiterID       *uuid.UUID     `gorm:"type:uuid;index" json:"waiter_id"`
	Waiter         *User          `gorm:"foreignKey:WaiterID" json:"waiter,omitempty"`
	CashierID      *uuid.UUID     `gorm:"type:uuid" json:"cashier_id"`
	OrderType      OrderType      `gorm:"default:'dine_in'" json:"order_type"`
	Status         OrderStatus    `gorm:"default:'open';index" json:"status"`
	Note           string         `json:"note"`
	Subtotal       float64        `gorm:"type:numeric(10,2);default:0" json:"subtotal"`
	TaxAmount      float64        `gorm:"type:numeric(10,2);default:0" json:"tax_amount"`
	ServiceCharge  float64        `gorm:"type:numeric(10,2);default:0" json:"service_charge"`
	DiscountAmount float64        `gorm:"type:numeric(10,2);default:0" json:"discount_amount"`
	Total          float64        `gorm:"type:numeric(10,2);default:0" json:"total"`
	PaymentMethod  string         `json:"payment_method"`
	PaymentStatus  PaymentStatus  `gorm:"default:'unpaid'" json:"payment_status"`
	PaidAt         *time.Time     `json:"paid_at"`
	LocalID        string         `gorm:"index" json:"local_id"`
	SyncedAt       *time.Time     `json:"synced_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}
