package services

import (
	"kimsha/internal/models"
	"kimsha/internal/repository"
	"kimsha/internal/ws"

	"github.com/google/uuid"
)

type KitchenService struct{ repo *repository.KitchenRepository }

func NewKitchenService(r *repository.KitchenRepository) *KitchenService {
	return &KitchenService{repo: r}
}

func (s *KitchenService) ActiveTickets(tenantID uuid.UUID) ([]models.KitchenTicket, error) {
	return s.repo.ActiveTickets(tenantID)
}

func (s *KitchenService) UpdateTicketStatus(id uuid.UUID, tenantID string, status models.TicketStatus) error {
	if err := s.repo.UpdateStatus(id, status); err != nil {
		return err
	}
	ticket, _ := s.repo.Find(id)
	ws.Default.Broadcast(tenantID, ws.Message{
		Type:    "ticket_updated",
		Payload: ticket,
	})
	return nil
}
