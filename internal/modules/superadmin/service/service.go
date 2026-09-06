package service

import (
	"context"
	"encoding/json"

	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"rentos-backend/internal/modules/superadmin/dto/request"
	"rentos-backend/internal/modules/superadmin/dto/response"
	"rentos-backend/internal/modules/superadmin/entity"
	"rentos-backend/internal/modules/superadmin/repository"

	// Auth module needed to create the first admin user for the tenant
	authentity "rentos-backend/internal/modules/auth/entity"
	authrepo "rentos-backend/internal/modules/auth/repository"
	"rentos-backend/pkg/password"
	pkgresp "rentos-backend/pkg/response"
	"rentos-backend/pkg/transaction"
)

// ---- Interfaces ----

type SuperAdminService interface {
	Login(ctx context.Context, req request.SuperAdminLogin) (*response.SuperAdminAuth, error)
	RegisterTenant(ctx context.Context, req request.RegisterTenant) (*response.RegisterTenantResult, error)
	ListTenants(ctx context.Context, filter request.ListTenantsFilter) ([]response.TenantListItem, error)
	GetTenant(ctx context.Context, tenantID string) (*response.TenantDetail, error)
	UpdateTenantStatus(ctx context.Context, tenantID string, req request.UpdateTenantStatus) error
	AssignSubscription(ctx context.Context, tenantID string, req request.AssignSubscription) error
	ListPlans(ctx context.Context) ([]response.SubscriptionPlan, error)
	GetStats(ctx context.Context) (*response.PlatformStats, error)
}

type accessTokenIssuer interface {
	IssueAccess(userID, tenantID, email string) (string, error)
}

// ---- Implementation ----

type superAdminService struct {
	db          *sqlx.DB
	tenants     repository.TenantRepository
	plans       repository.SubscriptionPlanRepository
	subs        repository.TenantSubscriptionRepository
	superAdmins repository.SuperAdminRepository
	users       authrepo.UserRepository
	jwtClient   accessTokenIssuer
}

func NewSuperAdminService(
	db *sqlx.DB,
	tenants repository.TenantRepository,
	plans repository.SubscriptionPlanRepository,
	subs repository.TenantSubscriptionRepository,
	superAdmins repository.SuperAdminRepository,
	users authrepo.UserRepository,
	jwtClient accessTokenIssuer,
) SuperAdminService {
	return &superAdminService{
		db: db, tenants: tenants, plans: plans,
		subs: subs, superAdmins: superAdmins,
		users: users, jwtClient: jwtClient,
	}
}

// Login authenticates a super admin and returns a JWT.
func (s *superAdminService) Login(ctx context.Context, req request.SuperAdminLogin) (*response.SuperAdminAuth, error) {
	admin, err := s.superAdmins.FindByEmail(ctx, s.db, req.Email)
	if err != nil {
		// Don't leak whether the email exists
		return nil, pkgresp.NewAppError(pkgresp.CodeUnauthorized, "email atau password salah")
	}

	if !admin.IsActive {
		return nil, pkgresp.NewAppError(pkgresp.CodeUnauthorized, "akun super admin dinonaktifkan")
	}

	if err := password.Verify(admin.PasswordHash, req.Password); err != nil {
		return nil, pkgresp.NewAppError(pkgresp.CodeUnauthorized, "email atau password salah")
	}

	// Issue JWT with special super_admin claim
	token, err := s.jwtClient.IssueAccess(admin.ID, "platform", admin.Email)
	if err != nil {
		return nil, fmt.Errorf("superAdminService.Login: issue token: %w", err)
	}

	_ = s.superAdmins.UpdateLastLogin(ctx, s.db, admin.ID)

	return &response.SuperAdminAuth{
		AccessToken: token,
		Admin: response.SuperAdminInfo{
			ID:       admin.ID,
			Email:    admin.Email,
			FullName: admin.FullName,
		},
	}, nil
}

// RegisterTenant creates a new tenant + admin user in a single transaction.
func (s *superAdminService) RegisterTenant(ctx context.Context, req request.RegisterTenant) (*response.RegisterTenantResult, error) {
	// Check slug uniqueness
	if exists, _ := s.tenants.SlugExists(ctx, s.db, req.BusinessSlug); exists {
		return nil, pkgresp.NewAppError(pkgresp.CodeConflict, "slug bisnis sudah digunakan, coba yang lain")
	}

	// Find the requested plan
	plan, err := s.plans.FindBySlug(ctx, s.db, req.Plan)
	if err != nil {
		return nil, pkgresp.NewAppError(pkgresp.CodeNotFound, "plan tidak ditemukan 2")
	}

	tenantID := uuid.NewString()
	adminID := uuid.NewString()

	hashPass, err := password.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("superAdminService.RegisterTenant: hash: %w", err)
	}

	now := time.Now().UTC()
	var trialEndsAt *time.Time
	if plan.Slug == "trial" {
		t := now.Add(14 * 24 * time.Hour)
		trialEndsAt = &t
	}

	err = transaction.WithTx(ctx, s.db, func(tx *sqlx.Tx) error {
		// 1. Create tenant
		tenant := &entity.Tenant{
			ID:                 tenantID,
			Name:               req.BusinessName,
			Slug:               req.BusinessSlug,
			Email:              req.Email,
			Phone:              req.BusinessPhone,
			Timezone:           "Asia/Jakarta",
			Locale:             "id-ID",
			Currency:           "IDR",
			SubscriptionStatus: entity.SubscriptionStatusTrial,
			TrialEndsAt:        trialEndsAt,
			Status:             "active",
		}
		if plan.Slug != "trial" {
			tenant.SubscriptionStatus = entity.SubscriptionStatusActive
			tenant.TrialEndsAt = nil
		}
		if err := s.tenants.Create(ctx, tx, tenant); err != nil {
			return fmt.Errorf("create tenant: %w", err)
		}

		// 2. Create admin user for the tenant
		nameParts := splitName(req.FullName)
		user := &authentity.User{
			ID:              adminID,
			TenantID:        tenantID,
			Username:        req.BusinessSlug + "_admin",
			Email:           req.Email,
			PasswordHash:    &hashPass,
			FirstName:       &nameParts[0],
			LastName:        &nameParts[1],
			EmailVerifiedAt: nil, // TODO: email verification flow
			IsActive:        true,
			Status:          "active",
		}
		if err := s.users.Create(ctx, tx, user); err != nil {
			return fmt.Errorf("create admin user: %w", err)
		}

		// 3. If paid plan, create subscription record
		if plan.Slug != "trial" {
			sub := &entity.TenantSubscription{
				ID:           uuid.NewString(),
				TenantID:     tenantID,
				PlanID:       plan.ID,
				BillingCycle: "monthly",
				PricePaid:    plan.PriceMonthly,
				StartsAt:     now,
				EndsAt:       now.Add(30 * 24 * time.Hour),
				AutoRenew:    true,
				Status:       "active",
			}
			if err := s.subs.Create(ctx, tx, sub); err != nil {
				return fmt.Errorf("create subscription: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &response.RegisterTenantResult{
		TenantID:   tenantID,
		TenantSlug: req.BusinessSlug,
		AdminEmail: req.Email,
		Plan:       plan.Name,
		Message:    "Registrasi berhasil! Silakan login untuk mulai menggunakan RentOS.",
	}, nil
}

func (s *superAdminService) ListTenants(ctx context.Context, filter request.ListTenantsFilter) ([]response.TenantListItem, error) {
	perPage, page := normPage(filter.PerPage, filter.Page)
	tenants, err := s.tenants.List(ctx, s.db, filter.Search, filter.SubscriptionStatus, filter.Status, perPage, (page-1)*perPage)
	if err != nil {
		return nil, err
	}
	out := make([]response.TenantListItem, 0, len(tenants))
	for _, t := range tenants {
		out = append(out, response.TenantListItem{
			ID: t.ID, Name: t.Name, Slug: t.Slug, Email: t.Email,
			Phone: t.Phone, SubscriptionStatus: t.SubscriptionStatus,
			TrialEndsAt: t.TrialEndsAt, Status: t.Status, CreatedAt: t.CreatedAt,
		})
	}
	return out, nil
}

func (s *superAdminService) GetTenant(ctx context.Context, tenantID string) (*response.TenantDetail, error) {
	t, err := s.tenants.FindByID(ctx, s.db, tenantID)
	if err != nil {
		return nil, pkgresp.NewAppError(pkgresp.CodeNotFound, "tenant tidak ditemukan")
	}
	return &response.TenantDetail{
		TenantListItem: response.TenantListItem{
			ID: t.ID, Name: t.Name, Slug: t.Slug, Email: t.Email,
			Phone: t.Phone, SubscriptionStatus: t.SubscriptionStatus,
			TrialEndsAt: t.TrialEndsAt, Status: t.Status, CreatedAt: t.CreatedAt,
		},
		Timezone: t.Timezone, Locale: t.Locale, Currency: t.Currency,
	}, nil
}

func (s *superAdminService) UpdateTenantStatus(ctx context.Context, tenantID string, req request.UpdateTenantStatus) error {
	return s.tenants.UpdateStatus(ctx, s.db, tenantID, req.Status)
}

func (s *superAdminService) AssignSubscription(ctx context.Context, tenantID string, req request.AssignSubscription) error {
	now := time.Now().UTC()
	sub := &entity.TenantSubscription{
		ID:           uuid.NewString(),
		TenantID:     tenantID,
		PlanID:       req.PlanID,
		BillingCycle: req.BillingCycle,
		PricePaid:    req.PricePaid,
		StartsAt:     now,
		EndsAt:       now.Add(time.Duration(req.DurationDays) * 24 * time.Hour),
		AutoRenew:    true,
		Status:       "active",
	}
	return transaction.WithTx(ctx, s.db, func(tx *sqlx.Tx) error {
		if err := s.subs.Create(ctx, tx, sub); err != nil {
			return err
		}
		return s.tenants.UpdateSubscriptionStatus(ctx, tx, tenantID, entity.SubscriptionStatusActive)
	})
}

func (s *superAdminService) ListPlans(ctx context.Context) ([]response.SubscriptionPlan, error) {
	plans, err := s.plans.ListActive(ctx, s.db)
	if err != nil {
		return nil, err
	}
	out := make([]response.SubscriptionPlan, 0, len(plans))
	for _, p := range plans {
		var features []string
		_ = json.Unmarshal(p.Features, &features)
		out = append(out, response.SubscriptionPlan{
			ID: p.ID, Name: p.Name, Slug: p.Slug, Description: p.Description,
			PriceMonthly: p.PriceMonthly, PriceYearly: p.PriceYearly,
			MaxAssets: p.MaxAssets, MaxUsers: p.MaxUsers, Features: features,
		})
	}
	return out, nil
}

func (s *superAdminService) GetStats(ctx context.Context) (*response.PlatformStats, error) {
	return s.tenants.GetPlatformStats(ctx, s.db)
}

// ---- helpers ----

func splitName(fullName string) [2]string {
	for i, ch := range fullName {
		if ch == ' ' {
			return [2]string{fullName[:i], fullName[i+1:]}
		}
	}
	return [2]string{fullName, ""}
}

func normPage(perPage, page int) (int, int) {
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	if page <= 0 {
		page = 1
	}
	return perPage, page
}
