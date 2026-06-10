package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name      string         `gorm:"not null" json:"name"`
	NameAm    string         `json:"name_am"`
	Icon      string         `json:"icon"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Items []MenuItem `gorm:"foreignKey:CategoryID" json:"items,omitempty"`
}
