package services

import (
	"encoding/json"
	"fmt"

	"kimsha/internal/models"
	"kimsha/internal/repository"
	"kimsha/internal/ws"
	"kimsha/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	menuRepo    *repository.MenuRepository
	kitchenRepo *repository.KitchenRepository
	tableRepo   *repository.TableRepository
}

func NewOrderService(
	or *repository.OrderRepository,
	mr *repository.MenuRepository,
	kr *repository.KitchenRepository,
	tr *repository.TableRepository,
) *OrderService {
	return &OrderService{orderRepo: or, menuRepo: mr, kitchenRepo: kr, tableRepo: tr}
}

type CreateOrderInput struct {
	TableID   *uuid.UUID        `json:"table_id"`
	WaiterID  *uuid.UUID        `json:"waiter_id"`
	OrderType models.OrderType  `json:"order_type"`
	Note      string            `json:"note"`
	LocalID   string            `json:"local_id"`
}

type AddItemInput struct {
	ItemID   uuid.UUID       `json:"item_id" validate:"required"`
	Quantity int             `json:"quantity" validate:"required,min=1"`
	Note     string          `json:"note"`
	AddOns   []AddOnSnapshot `json:"add_ons"`
}

type AddOnSnapshot struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (s *OrderService) List(tenantID uuid.UUID, status string, p pagination.Params) ([]models.Order, int64, error) {
	return s.orderRepo.List(tenantID, status, p)
}

func (s *OrderService) Create(tenantID uuid.UUID, in CreateOrderInput) (*models.Order, error) {
	order := &models.Order{
		ID:        uuid.New(),
		TenantID:  tenantID,
		TableID:   in.TableID,
		WaiterID:  in.WaiterID,
		OrderType: in.OrderType,
		Note:      in.Note,
		LocalID:   in.LocalID,
		Status:    models.OrderStatusOpen,
	}
	if order.OrderType == "" {
		order.OrderType = models.OrderTypeDineIn
	}
	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}
	if in.TableID != nil {
		_ = s.tableRepo.UpdateStatus(*in.TableID, tenantID, models.TableOccupied)
	}
	return order, nil
}

func (s *OrderService) Get(id, tenantID uuid.UUID) (*models.Order, error) {
	return s.orderRepo.Find(id, tenantID)
}

func (s *OrderService) UpdateStatus(id, tenantID uuid.UUID, status models.OrderStatus) error {
	return s.orderRepo.UpdateStatus(id, tenantID, status)
}

func (s *OrderService) AddItem(orderID, tenantID uuid.UUID, in AddItemInput) (*models.OrderItem, error) {
	menuItem, err := s.menuRepo.FindItem(in.ItemID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("item not found")
	}

	addOnsJSON, _ := json.Marshal(in.AddOns)
	item := &models.OrderItem{
		ID:            uuid.New(),
		OrderID:       orderID,
		ItemID:        in.ItemID,
		NameSnapshot:  menuItem.Name,
		PriceSnapshot: menuItem.Price,
		Quantity:      in.Quantity,
		Note:          in.Note,
		AddOns:        datatypes.JSON(addOnsJSON),
		Status:        models.OrderItemPending,
	}
	if err := s.orderRepo.AddItem(item); err != nil {
		return nil, err
	}

	// create kitchen ticket
	order, _ := s.orderRepo.Find(orderID, tenantID)
	tableRef := "Takeaway"
	if order != nil && order.Table != nil {
		tableRef = fmt.Sprintf("Table %d", order.Table.Number)
	}
	ticket := &models.KitchenTicket{
		ID:       uuid.New(),
		TenantID: tenantID,
		OrderID:  orderID,
		ItemID:   &item.ID,
		TableRef: tableRef,
		ItemName: menuItem.Name,
		Quantity: in.Quantity,
		Note:     in.Note,
		Status:   models.TicketQueued,
	}
	_ = s.kitchenRepo.Create(ticket)

	ws.Default.Broadcast(tenantID.String(), ws.Message{
		Type:    "new_ticket",
		Payload: ticket,
	})

	// recalculate order totals
	s.recalcOrder(orderID, tenantID)

	return item, nil
}

func (s *OrderService) UpdateItem(orderID, itemID, tenantID uuid.UUID, status models.OrderItemStatus) error {
	item, err := s.orderRepo.FindItem(itemID)
	if err != nil {
		return err
	}
	item.Status = status
	return s.orderRepo.UpdateItem(item)
}

func (s *OrderService) RemoveItem(orderID, itemID, tenantID uuid.UUID) error {
	return s.orderRepo.DeleteItem(itemID)
}

func (s *OrderService) recalcOrder(orderID, tenantID uuid.UUID) {
	order, err := s.orderRepo.Find(orderID, tenantID)
	if err != nil {
		return
	}
	var subtotal float64
	for _, item := range order.Items {
		subtotal += item.PriceSnapshot * float64(item.Quantity)
	}
	order.Subtotal = subtotal
	order.Total = subtotal + order.TaxAmount + order.ServiceCharge - order.DiscountAmount
	_ = s.orderRepo.Update(order)
}
