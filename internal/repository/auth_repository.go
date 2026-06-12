package repository

import (
	"kimsha/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository struct{ db *gorm.DB }

func NewAuthRepository(db *gorm.DB) *AuthRepository { return &AuthRepository{db: db} }

func (r *AuthRepository) CreateTenant(t *models.Tenant) error {
	return r.db.Create(t).Error
}

func (r *AuthRepository) FindTenantBySlug(slug string) (*models.Tenant, error) {
	var t models.Tenant
	err := r.db.Where("slug = ? AND is_active = true", slug).First(&t).Error
	return &t, err
}

func (r *AuthRepository) CreateUser(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *AuthRepository) FindUserByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.Where("email = ? AND is_active = true", email).First(&u).Error
	return &u, err
}

func (r *AuthRepository) FindUserByID(id uuid.UUID) (*models.User, error) {
	var u models.User
	err := r.db.Where("id = ? AND is_active = true", id).First(&u).Error
	return &u, err
}

func (r *AuthRepository) FindActiveUsersByTenant(tenantID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("tenant_id = ? AND is_active = true AND pin != ''", tenantID).Find(&users).Error
	return users, err
}

func (r *AuthRepository) UpdateLastLogin(userID uuid.UUID) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		Update("last_login", gorm.Expr("NOW()")).Error
}
