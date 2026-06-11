package repository

import (
	"kimsha/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminRepository struct{ db *gorm.DB }

func NewAdminRepository(db *gorm.DB) *AdminRepository { return &AdminRepository{db: db} }

type TenantListItem struct {
	models.Tenant
	OwnerName  string `json:"owner_name"`
	OwnerEmail string `json:"owner_email"`
	UserCount  int64  `json:"user_count"`
	OrderCount int64  `json:"order_count"`
}

type TenantDetail struct {
	models.Tenant
	Owner      *models.User  `json:"owner"`
	Users      []models.User `json:"users"`
	UserCount  int64         `json:"user_count"`
	OrderCount int64         `json:"order_count"`
	Revenue    float64       `json:"revenue"`
}

type PlatformStats struct {
	TotalTenants  int64   `json:"total_tenants"`
	ActiveTenants int64   `json:"active_tenants"`
	TotalUsers    int64   `json:"total_users"`
	TotalOrders   int64   `json:"total_orders"`
	TotalRevenue  float64 `json:"total_revenue"`
}

func (r *AdminRepository) FindSuperAdmin(email string) (*models.User, error) {
	var u models.User
	err := r.db.Where("email = ? AND role = ? AND is_active = true", email, models.RoleSuperAdmin).First(&u).Error
	return &u, err
}

func (r *AdminRepository) CreateSuperAdmin(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *AdminRepository) SuperAdminExists() (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("role = ?", models.RoleSuperAdmin).Count(&count).Error
	return count > 0, err
}

func (r *AdminRepository) ListTenants() ([]TenantListItem, error) {
	var tenants []models.Tenant
	if err := r.db.Order("created_at DESC").Find(&tenants).Error; err != nil {
		return nil, err
	}

	items := make([]TenantListItem, len(tenants))
	for i, t := range tenants {
		item := TenantListItem{Tenant: t}

		var owner models.User
		r.db.Where("tenant_id = ? AND role = ?", t.ID, models.RoleOwner).First(&owner)
		item.OwnerName = owner.Name
		item.OwnerEmail = owner.Email

		r.db.Model(&models.User{}).Where("tenant_id = ? AND deleted_at IS NULL", t.ID).Count(&item.UserCount)
		r.db.Model(&models.Order{}).Where("tenant_id = ?", t.ID).Count(&item.OrderCount)

		items[i] = item
	}
	return items, nil
}

func (r *AdminRepository) GetTenantDetail(id uuid.UUID) (*TenantDetail, error) {
	var tenant models.Tenant
	if err := r.db.First(&tenant, "id = ?", id).Error; err != nil {
		return nil, err
	}

	detail := &TenantDetail{Tenant: tenant}

	var owner models.User
	if err := r.db.Where("tenant_id = ? AND role = ?", id, models.RoleOwner).First(&owner).Error; err == nil {
		detail.Owner = &owner
	}

	r.db.Where("tenant_id = ? AND deleted_at IS NULL", id).Find(&detail.Users)
	r.db.Model(&models.User{}).Where("tenant_id = ? AND deleted_at IS NULL", id).Count(&detail.UserCount)
	r.db.Model(&models.Order{}).Where("tenant_id = ?", id).Count(&detail.OrderCount)
	r.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount),0)").Where("tenant_id = ?", id).Scan(&detail.Revenue)

	return detail, nil
}

func (r *AdminRepository) CreateTenantWithOwner(tenant *models.Tenant, owner *models.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(tenant).Error; err != nil {
			return err
		}
		owner.TenantID = &tenant.ID
		return tx.Create(owner).Error
	})
}

func (r *AdminRepository) SetTenantStatus(id uuid.UUID, active bool) error {
	return r.db.Model(&models.Tenant{}).Where("id = ?", id).Update("is_active", active).Error
}

func (r *AdminRepository) GetStats() (*PlatformStats, error) {
	stats := &PlatformStats{}
	r.db.Model(&models.Tenant{}).Count(&stats.TotalTenants)
	r.db.Model(&models.Tenant{}).Where("is_active = true").Count(&stats.ActiveTenants)
	r.db.Model(&models.User{}).Where("role != ? AND deleted_at IS NULL", models.RoleSuperAdmin).Count(&stats.TotalUsers)
	r.db.Model(&models.Order{}).Count(&stats.TotalOrders)
	r.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount),0)").Scan(&stats.TotalRevenue)
	return stats, nil
}
