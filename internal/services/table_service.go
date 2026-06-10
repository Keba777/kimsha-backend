package services

import (
	"kimsha/internal/models"
	"kimsha/internal/repository"

	"github.com/google/uuid"
)

type TableService struct{ repo *repository.TableRepository }

func NewTableService(r *repository.TableRepository) *TableService { return &TableService{repo: r} }

type TableInput struct {
	Number   int    `json:"number" validate:"required,min=1"`
	Name     string `json:"name"`
	NameAm   string `json:"name_am"`
	Capacity int    `json:"capacity"`
	Section  string `json:"section"`
}

func (s *TableService) List(tenantID uuid.UUID) ([]models.Table, error) {
	return s.repo.List(tenantID)
}

func (s *TableService) Create(tenantID uuid.UUID, in TableInput) (*models.Table, error) {
	t := &models.Table{
		ID:       uuid.New(),
		TenantID: tenantID,
		Number:   in.Number,
		Name:     in.Name,
		NameAm:   in.NameAm,
		Capacity: in.Capacity,
		Section:  in.Section,
		Status:   models.TableFree,
	}
	if t.Capacity == 0 {
		t.Capacity = 4
	}
	return t, s.repo.Create(t)
}

func (s *TableService) Update(id, tenantID uuid.UUID, in TableInput) (*models.Table, error) {
	t, err := s.repo.Find(id, tenantID)
	if err != nil {
		return nil, err
	}
	t.Number = in.Number
	t.Name = in.Name
	t.NameAm = in.NameAm
	t.Capacity = in.Capacity
	t.Section = in.Section
	return t, s.repo.Update(t)
}

func (s *TableService) Delete(id, tenantID uuid.UUID) error {
	return s.repo.Delete(id, tenantID)
}

func (s *TableService) UpdateStatus(id, tenantID uuid.UUID, status models.TableStatus) error {
	return s.repo.UpdateStatus(id, tenantID, status)
}
