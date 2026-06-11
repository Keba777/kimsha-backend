package services

import (
	"fmt"
	"strings"
	"time"

	"kimsha/internal/models"
	"kimsha/internal/repository"
	"kimsha/pkg/jwt"
	"kimsha/pkg/password"

	"github.com/google/uuid"
)

type AuthService struct {
	authRepo *repository.AuthRepository
	expiry   time.Duration
}

func NewAuthService(r *repository.AuthRepository, expiry time.Duration) *AuthService {
	return &AuthService{authRepo: r, expiry: expiry}
}

type RegisterInput struct {
	TenantName string `json:"tenant_name" validate:"required,min=2"`
	TenantSlug string `json:"tenant_slug" validate:"required,min=2,alphanum"`
	OwnerName  string `json:"owner_name" validate:"required,min=2"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type PINLoginInput struct {
	TenantSlug string `json:"tenant_slug" validate:"required"`
	PIN        string `json:"pin" validate:"required,len=4"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (s *AuthService) Register(input RegisterInput) (*AuthResponse, error) {
	tenant := &models.Tenant{
		ID:   uuid.New(),
		Name: input.TenantName,
		Slug: strings.ToLower(input.TenantSlug),
	}
	if err := s.authRepo.CreateTenant(tenant); err != nil {
		return nil, fmt.Errorf("tenant: %w", err)
	}

	hashed, err := password.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:       uuid.New(),
		TenantID: &tenant.ID,
		Name:     input.OwnerName,
		Email:    input.Email,
		Password: hashed,
		Role:     models.RoleOwner,
	}
	if err := s.authRepo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	token, err := jwt.Sign(user.ID, tenant.ID, string(user.Role), user.Name, s.expiry)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{Token: token, User: user}, nil
}

func (s *AuthService) Login(input LoginInput) (*AuthResponse, error) {
	user, err := s.authRepo.FindUserByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	if !password.Compare(user.Password, input.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}
	_ = s.authRepo.UpdateLastLogin(user.ID)

	tenantID := uuid.Nil
	if user.TenantID != nil {
		tenantID = *user.TenantID
	}
	token, err := jwt.Sign(user.ID, tenantID, string(user.Role), user.Name, s.expiry)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{Token: token, User: user}, nil
}

func (s *AuthService) PINLogin(input PINLoginInput) (*AuthResponse, error) {
	tenant, err := s.authRepo.FindTenantBySlug(input.TenantSlug)
	if err != nil {
		return nil, fmt.Errorf("tenant not found")
	}

	// PIN stored as bcrypt hash
	users, err := s.authRepo.FindUserByEmail("") // placeholder — use pin lookup
	_ = users
	_ = tenant
	_ = err
	// Real implementation: look up by tenant + compare PIN hash
	return nil, fmt.Errorf("pin login not yet implemented")
}
