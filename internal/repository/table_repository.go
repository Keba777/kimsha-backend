package repository

import (
	"kimsha/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TableRepository struct{ db *gorm.DB }

func NewTableRepository(db *gorm.DB) *TableRepository { return &TableRepository{db: db} }

func (r *TableRepository) List(tenantID uuid.UUID) ([]models.Table, error) {
	var tables []models.Table
	err := r.db.Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Order("number ASC").Find(&tables).Error
	return tables, err
}

func (r *TableRepository) Create(t *models.Table) error {
	return r.db.Create(t).Error
}

func (r *TableRepository) Find(id, tenantID uuid.UUID) (*models.Table, error) {
	var t models.Table
	err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&t).Error
	return &t, err
}

func (r *TableRepository) Update(t *models.Table) error {
	return r.db.Save(t).Error
}

func (r *TableRepository) Delete(id, tenantID uuid.UUID) error {
	return r.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.Table{}).Error
}

func (r *TableRepository) UpdateStatus(id, tenantID uuid.UUID, status models.TableStatus) error {
	return r.db.Model(&models.Table{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Update("status", status).Error
}
