package models

import (
	"time"

	"github.com/google/uuid"
)

type CashTransaction struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	UserID    *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type      string     `gorm:"not null" json:"type"` // open | close | in | out | sale | refund
	Amount    float64    `gorm:"type:numeric(10,2);not null" json:"amount"`
	Note      string     `json:"note"`
	Balance   float64    `gorm:"type:numeric(10,2)" json:"balance"`
	ShiftDate time.Time  `gorm:"type:date;default:current_date" json:"shift_date"`
	CreatedAt time.Time  `json:"created_at"`
}
