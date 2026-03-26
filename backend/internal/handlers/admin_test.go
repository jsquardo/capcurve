package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestValidAdminAuthorization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		headerValue string
		adminSecret string
		want        bool
	}{
		{
			name:        "missing header",
			headerValue: "",
			adminSecret: "super-secret",
			want:        false,
		},
		{
			name:        "missing configured secret",
			headerValue: "Bearer super-secret",
			adminSecret: "",
			want:        false,
		},
		{
			name:        "wrong auth scheme",
			headerValue: "Token super-secret",
			adminSecret: "super-secret",
			want:        false,
		},
		{
			name:        "wrong token",
			headerValue: "Bearer wrong-secret",
			adminSecret: "super-secret",
			want:        false,
		},
		{
			name:        "missing bearer token",
			headerValue: "Bearer   ",
			adminSecret: "super-secret",
			want:        false,
		},
		{
			name:        "matching token",
			headerValue: "Bearer super-secret",
			adminSecret: "super-secret",
			want:        true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, validAdminAuthorization(tt.headerValue, tt.adminSecret))
		})
	}
}

func TestAdminDashboardRejectsUnauthorizedRequests(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/dashboard", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &Handler{adminSecret: "super-secret"}

	err := handler.AdminDashboard(c)

	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
	require.JSONEq(t, `{"error":"unauthorized"}`, rec.Body.String())
}
