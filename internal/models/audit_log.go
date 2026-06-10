package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AuditLog struct {
	ID       uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	UserID   *uuid.UUID     `gorm:"type:uuid" json:"user_id"`
	Action   string         `gorm:"not null" json:"action"`
	Entity   string         `json:"entity"`
	EntityID string         `json:"entity_id"`
	OldValue datatypes.JSON `gorm:"type:jsonb" json:"old_value"`
	NewValue datatypes.JSON `gorm:"type:jsonb" json:"new_value"`
	IP       string         `json:"ip"`
	CreatedAt time.Time     `json:"created_at"`
}
