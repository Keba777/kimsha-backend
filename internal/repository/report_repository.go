package repository

import (
	"time"

	"kimsha/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportRepository struct{ db *gorm.DB }

func NewReportRepository(db *gorm.DB) *ReportRepository { return &ReportRepository{db: db} }

type TopItem struct {
	ItemID   string  `json:"item_id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Revenue  float64 `json:"revenue"`
}

type HourlySale struct {
	Hour    int     `json:"hour"`
	Revenue float64 `json:"revenue"`
	Orders  int     `json:"orders"`
}

type WaiterStat struct {
	WaiterID uuid.UUID `json:"waiter_id"`
	Name     string    `json:"name"`
	Orders   int       `json:"orders"`
	Revenue  float64   `json:"revenue"`
}

func (r *ReportRepository) TopItems(tenantID uuid.UUID, from, to time.Time, limit int) ([]TopItem, error) {
	var items []TopItem
	err := r.db.Raw(`
		SELECT oi.item_id::text, oi.name_snapshot as name,
			   SUM(oi.quantity) as quantity,
			   SUM(oi.quantity * oi.price_snapshot) as revenue
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		WHERE o.tenant_id = ? AND o.created_at >= ? AND o.created_at <= ?
		  AND o.status != ?
		GROUP BY oi.item_id, oi.name_snapshot
		ORDER BY quantity DESC
		LIMIT ?`, tenantID, from, to, models.OrderStatusCancelled, limit).
		Scan(&items).Error
	return items, err
}

func (r *ReportRepository) HourlySales(tenantID uuid.UUID, date time.Time) ([]HourlySale, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)
	var sales []HourlySale
	err := r.db.Raw(`
		SELECT EXTRACT(HOUR FROM created_at)::int as hour,
			   COALESCE(SUM(total),0) as revenue,
			   COUNT(*) as orders
		FROM orders
		WHERE tenant_id = ? AND created_at >= ? AND created_at < ? AND status != ?
		GROUP BY hour ORDER BY hour`, tenantID, start, end, models.OrderStatusCancelled).
		Scan(&sales).Error
	return sales, err
}

func (r *ReportRepository) WaiterStats(tenantID uuid.UUID, from, to time.Time) ([]WaiterStat, error) {
	var stats []WaiterStat
	err := r.db.Raw(`
		SELECT o.waiter_id, u.name,
			   COUNT(*) as orders,
			   COALESCE(SUM(o.total),0) as revenue
		FROM orders o
		JOIN users u ON u.id = o.waiter_id
		WHERE o.tenant_id = ? AND o.created_at >= ? AND o.created_at <= ? AND o.waiter_id IS NOT NULL
		GROUP BY o.waiter_id, u.name
		ORDER BY revenue DESC`, tenantID, from, to).
		Scan(&stats).Error
	return stats, err
}
