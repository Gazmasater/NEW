package server

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"project.com/internal/config"
	"project.com/internal/storage"
)

func TestPostandGetJSONGauge(t *testing.T) {
	serverCfg := config.InitServerConfig()

	tests := []struct {
		name   string
		path   string
		method string
		body   []byte
	}{
		{
			name:   "Update: Test POST request to /update/counter/test1/120",
			path:   "/update/counter/test1/120",
			method: "POST",
			body:   []byte(`{"delta": 120}`),
		},
		// {
		// 	name:   "Value: Test POST request to /value/counter/test1",
		// 	path:   "/value/counter/test1",
		// 	method: "POST",
		// 	body:   nil,
		// },

		// {
		// 	name:   "Update: Test POST request to /update/counter/test2/120",
		// 	path:   "/update/counter/test2/120",
		// 	method: "POST",
		// 	body:   []byte(`{"delta": 120}`),
		// },
		// {
		// 	name:   "Value: Test POST request to /value/counter/test2",
		// 	path:   "/value/counter/test2",
		// 	method: "POST",
		// 	body:   nil,
		// },
	}

	mockStorage := storage.NewMemStorage()
	mockLogger := zap.NewExample()
	mockConfig := &config.ServerConfig{}

	db, err := sql.Open("postgres", serverCfg.DatabaseDSN)
	if err != nil {
		t.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	fieldData := app{
		Storage: mockStorage,
		Logger:  mockLogger,
		Config:  mockConfig,
		DB:      db,
	}

	mc := &app{
		Storage: fieldData.Storage,
		Logger:  fieldData.Logger,
		Config:  fieldData.Config,
		DB:      fieldData.DB,
	}

	r := chi.NewRouter()
	r.Post("/update/", mc.updateHandlerJSON)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(tt.body))
			rec := httptest.NewRecorder()

			mc.updateHandlerJSON(rec, req)

			responseBody := rec.Body.String()
			t.Logf("Response body for test '%s': %s", tt.name, responseBody)

			if rec.Code != http.StatusOK {
				t.Errorf("Expected status OK, got: %d", rec.Code)
			}

		})

	}
}
