package repository

import (
	"kimsha/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SyncRepository struct{ db *gorm.DB }

func NewSyncRepository(db *gorm.DB) *SyncRepository { return &SyncRepository{db: db} }

func (r *SyncRepository) CreateLog(log *models.SyncLog) error {
	return r.db.Create(log).Error
}

func (r *SyncRepository) UpdateLog(id uuid.UUID, serverID *uuid.UUID, status, errMsg string) error {
	updates := map[string]interface{}{
		"status": status,
		"error":  errMsg,
	}
	if serverID != nil {
		updates["server_id"] = *serverID
		updates["synced_at"] = gorm.Expr("NOW()")
	}
	return r.db.Model(&models.SyncLog{}).Where("id = ?", id).Updates(updates).Error
}
