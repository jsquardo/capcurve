package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jsquardo/capcurve/internal/syncjob"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestAdminDashboardPreflightIncludesCORSHeaders(t *testing.T) {
	t.Parallel()

	e := echo.New()
	RegisterRoutes(e, nil, syncjob.NewStatusStore(false), "super-secret")

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/admin/dashboard", nil)
	req.Header.Set(echo.HeaderOrigin, "http://localhost:5173")
	req.Header.Set(echo.HeaderAccessControlRequestMethod, http.MethodGet)
	req.Header.Set(echo.HeaderAccessControlRequestHeaders, echo.HeaderAuthorization)

	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Equal(t, "http://localhost:5173", rec.Header().Get(echo.HeaderAccessControlAllowOrigin))
	require.Contains(t, rec.Header().Get(echo.HeaderAccessControlAllowMethods), http.MethodGet)
	require.Contains(t, rec.Header().Get(echo.HeaderAccessControlAllowHeaders), echo.HeaderAuthorization)
}
