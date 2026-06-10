package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OrderItemStatus string

const (
	OrderItemPending   OrderItemStatus = "pending"
	OrderItemCooking   OrderItemStatus = "cooking"
	OrderItemReady     OrderItemStatus = "ready"
	OrderItemServed    OrderItemStatus = "served"
	OrderItemCancelled OrderItemStatus = "cancelled"
)

type OrderItem struct {
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID       uuid.UUID       `gorm:"type:uuid;not null;index" json:"order_id"`
	ItemID        uuid.UUID       `gorm:"type:uuid;not null" json:"item_id"`
	MenuItem      *MenuItem       `gorm:"foreignKey:ItemID" json:"menu_item,omitempty"`
	NameSnapshot  string          `gorm:"not null" json:"name_snapshot"`
	PriceSnapshot float64         `gorm:"type:numeric(10,2);not null" json:"price_snapshot"`
	Quantity      int             `gorm:"default:1" json:"quantity"`
	Note          string          `json:"note"`
	Status        OrderItemStatus `gorm:"default:'pending'" json:"status"`
	AddOns        datatypes.JSON  `gorm:"type:jsonb;default:'[]'" json:"add_ons"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}
