package repository

import (
	"kimsha/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KitchenRepository struct{ db *gorm.DB }

func NewKitchenRepository(db *gorm.DB) *KitchenRepository { return &KitchenRepository{db: db} }

func (r *KitchenRepository) ActiveTickets(tenantID uuid.UUID) ([]models.KitchenTicket, error) {
	var tickets []models.KitchenTicket
	err := r.db.Where("tenant_id = ? AND status IN ?", tenantID,
		[]string{string(models.TicketQueued), string(models.TicketCooking)}).
		Order("priority DESC, created_at ASC").Find(&tickets).Error
	return tickets, err
}

func (r *KitchenRepository) Create(t *models.KitchenTicket) error {
	return r.db.Create(t).Error
}

func (r *KitchenRepository) Find(id uuid.UUID) (*models.KitchenTicket, error) {
	var t models.KitchenTicket
	err := r.db.First(&t, "id = ?", id).Error
	return &t, err
}

func (r *KitchenRepository) UpdateStatus(id uuid.UUID, status models.TicketStatus) error {
	updates := map[string]interface{}{"status": status}
	if status == models.TicketCooking {
		updates["started_at"] = gorm.Expr("NOW()")
	}
	if status == models.TicketDone {
		updates["done_at"] = gorm.Expr("NOW()")
	}
	return r.db.Model(&models.KitchenTicket{}).Where("id = ?", id).Updates(updates).Error
}
