package services

import (
	"kimsha/internal/models"
	"kimsha/internal/repository"

	"github.com/google/uuid"
)

type MenuService struct{ repo *repository.MenuRepository }

func NewMenuService(r *repository.MenuRepository) *MenuService { return &MenuService{repo: r} }

// Categories

func (s *MenuService) ListCategories(tenantID uuid.UUID) ([]models.Category, error) {
	return s.repo.ListCategories(tenantID)
}

type CategoryInput struct {
	Name      string `json:"name" validate:"required"`
	NameAm    string `json:"name_am"`
	Icon      string `json:"icon"`
	SortOrder int    `json:"sort_order"`
}

func (s *MenuService) CreateCategory(tenantID uuid.UUID, in CategoryInput) (*models.Category, error) {
	c := &models.Category{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Name:      in.Name,
		NameAm:    in.NameAm,
		Icon:      in.Icon,
		SortOrder: in.SortOrder,
	}
	return c, s.repo.CreateCategory(c)
}

func (s *MenuService) UpdateCategory(id, tenantID uuid.UUID, in CategoryInput) (*models.Category, error) {
	c, err := s.repo.FindCategory(id, tenantID)
	if err != nil {
		return nil, err
	}
	c.Name = in.Name
	c.NameAm = in.NameAm
	c.Icon = in.Icon
	c.SortOrder = in.SortOrder
	return c, s.repo.UpdateCategory(c)
}

func (s *MenuService) DeleteCategory(id, tenantID uuid.UUID) error {
	return s.repo.DeleteCategory(id, tenantID)
}

// Items

type ItemInput struct {
	CategoryID    *uuid.UUID       `json:"category_id"`
	Name          string           `json:"name" validate:"required"`
	NameAm        string           `json:"name_am"`
	Description   string           `json:"description"`
	DescriptionAm string           `json:"description_am"`
	Price         float64          `json:"price" validate:"required,gt=0"`
	ImageURL      string           `json:"image_url"`
	ItemType      models.ItemType  `json:"item_type"`
	IsAvailable   bool             `json:"is_available"`
	IsFeatured    bool             `json:"is_featured"`
	PrepTimeMin   int              `json:"prep_time_min"`
	SortOrder     int              `json:"sort_order"`
}

func (s *MenuService) ListItems(tenantID uuid.UUID, categoryID *uuid.UUID, availableOnly bool) ([]models.MenuItem, error) {
	return s.repo.ListItems(tenantID, categoryID, availableOnly)
}

func (s *MenuService) CreateItem(tenantID uuid.UUID, in ItemInput) (*models.MenuItem, error) {
	item := &models.MenuItem{
		ID:            uuid.New(),
		TenantID:      tenantID,
		CategoryID:    in.CategoryID,
		Name:          in.Name,
		NameAm:        in.NameAm,
		Description:   in.Description,
		DescriptionAm: in.DescriptionAm,
		Price:         in.Price,
		ImageURL:      in.ImageURL,
		ItemType:      in.ItemType,
		IsAvailable:   in.IsAvailable,
		IsFeatured:    in.IsFeatured,
		PrepTimeMin:   in.PrepTimeMin,
		SortOrder:     in.SortOrder,
	}
	if item.ItemType == "" {
		item.ItemType = models.ItemTypeFood
	}
	return item, s.repo.CreateItem(item)
}

func (s *MenuService) UpdateItem(id, tenantID uuid.UUID, in ItemInput) (*models.MenuItem, error) {
	item, err := s.repo.FindItem(id, tenantID)
	if err != nil {
		return nil, err
	}
	item.CategoryID = in.CategoryID
	item.Name = in.Name
	item.NameAm = in.NameAm
	item.Description = in.Description
	item.DescriptionAm = in.DescriptionAm
	item.Price = in.Price
	item.ImageURL = in.ImageURL
	item.ItemType = in.ItemType
	item.IsAvailable = in.IsAvailable
	item.IsFeatured = in.IsFeatured
	item.PrepTimeMin = in.PrepTimeMin
	item.SortOrder = in.SortOrder
	return item, s.repo.UpdateItem(item)
}

func (s *MenuService) ToggleAvailability(id, tenantID uuid.UUID) (*models.MenuItem, error) {
	item, err := s.repo.FindItem(id, tenantID)
	if err != nil {
		return nil, err
	}
	item.IsAvailable = !item.IsAvailable
	return item, s.repo.UpdateItem(item)
}

func (s *MenuService) DeleteItem(id, tenantID uuid.UUID) error {
	return s.repo.DeleteItem(id, tenantID)
}

func (s *MenuService) GetItem(id, tenantID uuid.UUID) (*models.MenuItem, error) {
	return s.repo.FindItem(id, tenantID)
}
