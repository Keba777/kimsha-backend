package repository

import (
	"kimsha/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) List(tenantID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Order("created_at ASC").Find(&users).Error
	return users, err
}

func (r *UserRepository) Find(id, tenantID uuid.UUID) (*models.User, error) {
	var u models.User
	err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&u).Error
	return &u, err
}

func (r *UserRepository) Create(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepository) Update(u *models.User) error {
	return r.db.Save(u).Error
}

func (r *UserRepository) Delete(id, tenantID uuid.UUID) error {
	return r.db.Model(&models.User{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Update("is_active", false).Error
}

func (r *UserRepository) FindByPIN(tenantID uuid.UUID, pin string) (*models.User, error) {
	var u models.User
	err := r.db.Where("tenant_id = ? AND pin = ? AND is_active = true", tenantID, pin).First(&u).Error
	return &u, err
}
