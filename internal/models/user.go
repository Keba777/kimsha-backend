package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleSuperAdmin UserRole = "super_admin"
	RoleOwner      UserRole = "owner"
	RoleManager    UserRole = "manager"
	RoleWaiter     UserRole = "waiter"
	RoleKitchen    UserRole = "kitchen"
	RoleCashier    UserRole = "cashier"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID  *uuid.UUID     `gorm:"type:uuid;index" json:"tenant_id"`
	Tenant    Tenant         `gorm:"foreignKey:TenantID" json:"-"`
	Name      string         `gorm:"not null" json:"name"`
	NameAm    string         `json:"name_am"`
	Email     string         `gorm:"uniqueIndex" json:"email"`
	Phone     string         `json:"phone"`
	Password  string         `gorm:"not null" json:"-"`
	Role      UserRole       `gorm:"not null" json:"role"`
	PIN       string         `json:"-"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	LastLogin *time.Time     `json:"last_login"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
