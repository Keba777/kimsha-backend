package services

import (
	"time"

	"kimsha/internal/repository"

	"github.com/google/uuid"
)

type ReportService struct {
	reportRepo *repository.ReportRepository
	orderRepo  *repository.OrderRepository
}

func NewReportService(rr *repository.ReportRepository, or *repository.OrderRepository) *ReportService {
	return &ReportService{reportRepo: rr, orderRepo: or}
}

func (s *ReportService) Daily(tenantID uuid.UUID, date time.Time) (map[string]interface{}, error) {
	return s.orderRepo.DailySummary(tenantID, date)
}

func (s *ReportService) TopItems(tenantID uuid.UUID, from, to time.Time) ([]repository.TopItem, error) {
	return s.reportRepo.TopItems(tenantID, from, to, 10)
}

func (s *ReportService) HourlySales(tenantID uuid.UUID, date time.Time) ([]repository.HourlySale, error) {
	return s.reportRepo.HourlySales(tenantID, date)
}

func (s *ReportService) WaiterStats(tenantID uuid.UUID, from, to time.Time) ([]repository.WaiterStat, error) {
	return s.reportRepo.WaiterStats(tenantID, from, to)
}
