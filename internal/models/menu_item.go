package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemType string

const (
	ItemTypeFood  ItemType = "food"
	ItemTypeDrink ItemType = "drink"
	ItemTypeOther ItemType = "other"
)

type MenuItem struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	CategoryID    *uuid.UUID     `gorm:"type:uuid;index" json:"category_id"`
	Category      *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name          string         `gorm:"not null" json:"name"`
	NameAm        string         `json:"name_am"`
	Description   string         `json:"description"`
	DescriptionAm string         `json:"description_am"`
	Price         float64        `gorm:"type:numeric(10,2);not null" json:"price"`
	ImageURL      string         `json:"image_url"`
	ItemType      ItemType       `gorm:"default:'food'" json:"item_type"`
	IsAvailable   bool           `gorm:"default:true" json:"is_available"`
	IsFeatured    bool           `gorm:"default:false" json:"is_featured"`
	PrepTimeMin   int            `gorm:"default:10" json:"prep_time_min"`
	SortOrder     int            `gorm:"default:0" json:"sort_order"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	AddOns []AddOn `gorm:"foreignKey:ItemID" json:"add_ons,omitempty"`
}
