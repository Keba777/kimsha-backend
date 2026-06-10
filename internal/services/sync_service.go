package services

import (
	"encoding/json"
	"fmt"
	"time"

	"kimsha/internal/models"
	"kimsha/internal/repository"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type SyncService struct {
	syncRepo  *repository.SyncRepository
	orderRepo *repository.OrderRepository
	menuRepo  *repository.MenuRepository
	tableRepo *repository.TableRepository
}

func NewSyncService(
	sr *repository.SyncRepository,
	or *repository.OrderRepository,
	mr *repository.MenuRepository,
	tr *repository.TableRepository,
) *SyncService {
	return &SyncService{syncRepo: sr, orderRepo: or, menuRepo: mr, tableRepo: tr}
}

type SyncOperation struct {
	LocalID  string          `json:"local_id"`
	Entity   string          `json:"entity"`
	Action   string          `json:"action"` // create | update
	Payload  json.RawMessage `json:"payload"`
}

type PushRequest struct {
	DeviceID   string          `json:"device_id"`
	TenantID   uuid.UUID       `json:"tenant_id"`
	Operations []SyncOperation `json:"operations"`
}

type PushResult struct {
	LocalID  string `json:"local_id"`
	ServerID string `json:"server_id"`
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
}

func (s *SyncService) Push(tenantID uuid.UUID, req PushRequest) ([]PushResult, error) {
	results := make([]PushResult, 0, len(req.Operations))
	for _, op := range req.Operations {
		result := PushResult{LocalID: op.LocalID}
		serverID, err := s.applyOperation(tenantID, op)
		if err != nil {
			result.Status = "failed"
			result.Error = err.Error()
		} else {
			result.Status = "synced"
			result.ServerID = serverID.String()
		}

		log := &models.SyncLog{
			ID:       uuid.New(),
			TenantID: tenantID,
			DeviceID: req.DeviceID,
			Entity:   op.Entity,
			LocalID:  op.LocalID,
			Action:   op.Action,
			Payload:  datatypes.JSON(op.Payload),
			Status:   result.Status,
			Error:    result.Error,
		}
		if result.Status == "synced" {
			sid, _ := uuid.Parse(result.ServerID)
			log.ServerID = &sid
			now := time.Now()
			log.SyncedAt = &now
		}
		_ = s.syncRepo.CreateLog(log)
		results = append(results, result)
	}
	return results, nil
}

func (s *SyncService) applyOperation(tenantID uuid.UUID, op SyncOperation) (uuid.UUID, error) {
	switch op.Entity {
	case "orders":
		return s.applyOrder(tenantID, op)
	default:
		return uuid.Nil, fmt.Errorf("unknown entity: %s", op.Entity)
	}
}

func (s *SyncService) applyOrder(tenantID uuid.UUID, op SyncOperation) (uuid.UUID, error) {
	var input CreateOrderInput
	if err := json.Unmarshal(op.Payload, &input); err != nil {
		return uuid.Nil, err
	}
	order, err := (&OrderService{orderRepo: s.orderRepo, tableRepo: s.tableRepo}).
		Create(tenantID, input)
	if err != nil {
		return uuid.Nil, err
	}
	return order.ID, nil
}

func (s *SyncService) MenuSnapshot(tenantID uuid.UUID) (map[string]interface{}, error) {
	cats, err := s.menuRepo.ListCategories(tenantID)
	if err != nil {
		return nil, err
	}
	items, err := s.menuRepo.ListItems(tenantID, nil, false)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"categories": cats,
		"items":      items,
		"synced_at":  time.Now(),
	}, nil
}

func (s *SyncService) TableSnapshot(tenantID uuid.UUID) ([]models.Table, error) {
	return s.tableRepo.List(tenantID)
}
