package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name          string         `gorm:"not null" json:"name"`
	NameAm        string         `json:"name_am"`
	Slug          string         `gorm:"uniqueIndex;not null" json:"slug"`
	Phone         string         `json:"phone"`
	Address       string         `json:"address"`
	Timezone      string         `gorm:"default:'Africa/Addis_Ababa'" json:"timezone"`
	Currency      string         `gorm:"default:'ETB'" json:"currency"`
	TaxRate       float64        `gorm:"type:numeric(5,2);default:0" json:"tax_rate"`
	ServiceCharge float64        `gorm:"type:numeric(5,2);default:0" json:"service_charge"`
	Plan          string         `gorm:"default:'free'" json:"plan"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
