package repository

import (
	"strings"
	"time"

	"kimsha/internal/models"
	"kimsha/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct{ db *gorm.DB }

func NewOrderRepository(db *gorm.DB) *OrderRepository { return &OrderRepository{db: db} }

func (r *OrderRepository) List(tenantID uuid.UUID, status string, p pagination.Params) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64
	q := r.db.Where("tenant_id = ? AND deleted_at IS NULL", tenantID)
	if status != "" {
		parts := strings.Split(status, ",")
		if len(parts) == 1 {
			q = q.Where("status = ?", parts[0])
		} else {
			q = q.Where("status IN ?", parts)
		}
	}
	q.Model(&models.Order{}).Count(&total)
	err := q.Preload("Table").Preload("Waiter").Preload("Items").
		Order("created_at DESC").Limit(p.Limit).Offset(p.Offset()).Find(&orders).Error
	return orders, total, err
}

func (r *OrderRepository) Create(o *models.Order) error {
	return r.db.Create(o).Error
}

func (r *OrderRepository) Find(id, tenantID uuid.UUID) (*models.Order, error) {
	var o models.Order
	err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).
		Preload("Table").Preload("Waiter").Preload("Items.MenuItem").First(&o).Error
	return &o, err
}

func (r *OrderRepository) Update(o *models.Order) error {
	return r.db.Save(o).Error
}

func (r *OrderRepository) UpdateStatus(id, tenantID uuid.UUID, status models.OrderStatus) error {
	return r.db.Model(&models.Order{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Update("status", status).Error
}

func (r *OrderRepository) AddItem(item *models.OrderItem) error {
	return r.db.Create(item).Error
}

func (r *OrderRepository) FindItem(itemID uuid.UUID) (*models.OrderItem, error) {
	var item models.OrderItem
	err := r.db.First(&item, "id = ?", itemID).Error
	return &item, err
}

func (r *OrderRepository) UpdateItem(item *models.OrderItem) error {
	return r.db.Save(item).Error
}

func (r *OrderRepository) DeleteItem(itemID uuid.UUID) error {
	return r.db.Delete(&models.OrderItem{}, "id = ?", itemID).Error
}

func (r *OrderRepository) DailySummary(tenantID uuid.UUID, date time.Time) (map[string]interface{}, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	var result struct {
		TotalOrders  int64   `json:"total_orders"`
		TotalRevenue float64 `json:"total_revenue"`
		PaidOrders   int64   `json:"paid_orders"`
	}
	err := r.db.Model(&models.Order{}).
		Where("tenant_id = ? AND created_at >= ? AND created_at < ? AND status != ?",
			tenantID, start, end, models.OrderStatusCancelled).
		Select("COUNT(*) as total_orders, COALESCE(SUM(total),0) as total_revenue, COUNT(CASE WHEN payment_status='paid' THEN 1 END) as paid_orders").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"total_orders":  result.TotalOrders,
		"total_revenue": result.TotalRevenue,
		"paid_orders":   result.PaidOrders,
	}, nil
}
