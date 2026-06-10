package models

import (
	"time"

	"github.com/google/uuid"
)

type AddOn struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	ItemID     *uuid.UUID `gorm:"type:uuid;index" json:"item_id"`
	Name       string     `gorm:"not null" json:"name"`
	NameAm     string     `json:"name_am"`
	Price      float64    `gorm:"type:numeric(10,2);default:0" json:"price"`
	IsRequired bool       `gorm:"default:false" json:"is_required"`
	MaxSelect  int        `gorm:"default:1" json:"max_select"`
	CreatedAt  time.Time  `json:"created_at"`
}
