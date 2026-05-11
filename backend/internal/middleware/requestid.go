package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/macunaimaa/dashboard/backend/internal/service"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), service.CtxKeyRequestID, id)
		ctx = context.WithValue(ctx, service.CtxKeyIP, r.RemoteAddr)
		w.Header().Set("X-Request-ID", id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
