package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type SyncLog struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	DeviceID  string         `gorm:"not null;index" json:"device_id"`
	Entity    string         `gorm:"not null" json:"entity"`
	LocalID   string         `gorm:"not null" json:"local_id"`
	ServerID  *uuid.UUID     `gorm:"type:uuid" json:"server_id"`
	Action    string         `gorm:"not null" json:"action"` // create | update | delete
	Payload   datatypes.JSON `gorm:"type:jsonb" json:"payload"`
	Status    string         `gorm:"default:'pending'" json:"status"` // pending | synced | conflict | failed
	Error     string         `json:"error"`
	SyncedAt  *time.Time     `json:"synced_at"`
	CreatedAt time.Time      `json:"created_at"`
}
