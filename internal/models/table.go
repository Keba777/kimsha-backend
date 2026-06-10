package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TableStatus string

const (
	TableFree     TableStatus = "free"
	TableOccupied TableStatus = "occupied"
	TableReserved TableStatus = "reserved"
	TableCleaning TableStatus = "cleaning"
)

type Table struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Number    int            `gorm:"not null" json:"number"`
	Name      string         `json:"name"`
	NameAm    string         `json:"name_am"`
	Capacity  int            `gorm:"default:4" json:"capacity"`
	Section   string         `json:"section"`
	Status    TableStatus    `gorm:"default:'free'" json:"status"`
	QRCode    string         `json:"qr_code"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
