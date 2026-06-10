package models

import (
	"time"

	"github.com/google/uuid"
)

type TicketStatus string

const (
	TicketQueued    TicketStatus = "queued"
	TicketCooking   TicketStatus = "cooking"
	TicketDone      TicketStatus = "done"
	TicketCancelled TicketStatus = "cancelled"
)

type KitchenTicket struct {
	ID        uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID  uuid.UUID    `gorm:"type:uuid;not null;index" json:"tenant_id"`
	OrderID   uuid.UUID    `gorm:"type:uuid;not null;index" json:"order_id"`
	Order     *Order       `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	ItemID    *uuid.UUID   `gorm:"type:uuid" json:"item_id"`
	TableRef  string       `json:"table_ref"`
	ItemName  string       `gorm:"not null" json:"item_name"`
	Quantity  int          `gorm:"not null" json:"quantity"`
	Note      string       `json:"note"`
	Priority  int          `gorm:"default:0" json:"priority"`
	Status    TicketStatus `gorm:"default:'queued';index" json:"status"`
	StartedAt *time.Time   `json:"started_at"`
	DoneAt    *time.Time   `json:"done_at"`
	CreatedAt time.Time    `json:"created_at"`
}
