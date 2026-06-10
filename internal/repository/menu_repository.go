package repository

import (
	"kimsha/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MenuRepository struct{ db *gorm.DB }

func NewMenuRepository(db *gorm.DB) *MenuRepository { return &MenuRepository{db: db} }

// Categories

func (r *MenuRepository) ListCategories(tenantID uuid.UUID) ([]models.Category, error) {
	var cats []models.Category
	err := r.db.Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Order("sort_order ASC, created_at ASC").Find(&cats).Error
	return cats, err
}

func (r *MenuRepository) CreateCategory(c *models.Category) error {
	return r.db.Create(c).Error
}

func (r *MenuRepository) FindCategory(id, tenantID uuid.UUID) (*models.Category, error) {
	var c models.Category
	err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&c).Error
	return &c, err
}

func (r *MenuRepository) UpdateCategory(c *models.Category) error {
	return r.db.Save(c).Error
}

func (r *MenuRepository) DeleteCategory(id, tenantID uuid.UUID) error {
	return r.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.Category{}).Error
}

// Items

func (r *MenuRepository) ListItems(tenantID uuid.UUID, categoryID *uuid.UUID, availableOnly bool) ([]models.MenuItem, error) {
	var items []models.MenuItem
	q := r.db.Where("tenant_id = ? AND deleted_at IS NULL", tenantID)
	if categoryID != nil {
		q = q.Where("category_id = ?", *categoryID)
	}
	if availableOnly {
		q = q.Where("is_available = true")
	}
	err := q.Preload("Category").Preload("AddOns").
		Order("sort_order ASC, created_at ASC").Find(&items).Error
	return items, err
}

func (r *MenuRepository) CreateItem(item *models.MenuItem) error {
	return r.db.Create(item).Error
}

func (r *MenuRepository) FindItem(id, tenantID uuid.UUID) (*models.MenuItem, error) {
	var item models.MenuItem
	err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).
		Preload("Category").Preload("AddOns").First(&item).Error
	return &item, err
}

func (r *MenuRepository) UpdateItem(item *models.MenuItem) error {
	return r.db.Save(item).Error
}

func (r *MenuRepository) DeleteItem(id, tenantID uuid.UUID) error {
	return r.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.MenuItem{}).Error
}

// AddOns

func (r *MenuRepository) CreateAddOn(a *models.AddOn) error {
	return r.db.Create(a).Error
}

func (r *MenuRepository) FindAddOn(id uuid.UUID) (*models.AddOn, error) {
	var a models.AddOn
	err := r.db.First(&a, "id = ?", id).Error
	return &a, err
}

func (r *MenuRepository) UpdateAddOn(a *models.AddOn) error {
	return r.db.Save(a).Error
}

func (r *MenuRepository) DeleteAddOn(id uuid.UUID) error {
	return r.db.Delete(&models.AddOn{}, "id = ?", id).Error
}
