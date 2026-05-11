package service

import (
	"context"
	"time"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
	"github.com/macunaimaa/dashboard/backend/internal/repository"
)

type AuditService struct {
	repo repository.AuditRepository
}

func NewAuditService(repo repository.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) Log(ctx context.Context, tenantID, userID, action, resource string, detail interface{}) {
	entry := &domain.AuditEntry{
		TenantID:  tenantID,
		UserID:    userID,
		Action:    action,
		Resource:  resource,
		Detail:    detail,
		IP:        IPFromCtx(ctx),
		RequestID: RequestIDFromCtx(ctx),
		Timestamp: time.Now(),
	}
	// Fire and forget — audit logging should not block the request
	go func() {
		_ = s.repo.Insert(context.Background(), entry)
	}()
}

// Context key helpers — these are set by middleware

type ctxKey string

const (
	CtxKeyTenantID  ctxKey = "tenant_id"
	CtxKeyUserID    ctxKey = "user_id"
	CtxKeyRole      ctxKey = "role"
	CtxKeyRequestID ctxKey = "request_id"
	CtxKeyIP        ctxKey = "ip"
)

func TenantIDFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(CtxKeyTenantID).(string); ok {
		return v
	}
	return ""
}

func UserIDFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(CtxKeyUserID).(string); ok {
		return v
	}
	return ""
}

func RoleFromCtx(ctx context.Context) domain.Role {
	if v, ok := ctx.Value(CtxKeyRole).(domain.Role); ok {
		return v
	}
	return ""
}

func RequestIDFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(CtxKeyRequestID).(string); ok {
		return v
	}
	return ""
}

func IPFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(CtxKeyIP).(string); ok {
		return v
	}
	return ""
}
