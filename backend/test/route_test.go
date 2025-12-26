package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-expense-management-system/internal/delivery/http/route"
	"go-expense-management-system/internal/messages"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api")

	config := route.RouteConfig{
		Router: router,
	}

	config.RegisterApiRoutes(api)
	config.RegisterPublicRoutes()

	return router
}

func TestHealthEndpoints(t *testing.T) {
	router := setupRouter()

	tests := []string{"/health", "/api/health"}
	for _, path := range tests {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code, path)

		var payload map[string]string
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &payload), path)
		require.Equal(t, "ok", payload["status"], path)
	}
}

func TestWelcomeEndpoints(t *testing.T) {
	router := setupRouter()

	tests := map[string]string{
		"/":     messages.WelcomeMessage,
		"/api":  messages.WelcomeMessage,
	}

	for path, expected := range tests {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code, path)

		var payload map[string]string
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &payload), path)
		require.Equal(t, expected, payload["message"], path)
	}
}

func TestMetricsEndpoint(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/metrics", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var payload map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &payload))
	require.Equal(t, "ok", payload["status"])
	require.NotEmpty(t, payload["uptime"])
}
