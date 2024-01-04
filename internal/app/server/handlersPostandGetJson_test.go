package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"project.com/internal/config"
	"project.com/internal/storage"
)

type testFields struct {
	Storage *storage.MemStorage
	Logger  *zap.Logger
	Config  *config.ServerConfig
	// Добавьте вашу базу данных, если необходимо
}

func TestUpdateHandlerJSON(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		method string
		body   []byte
	}{
		{
			name:   "Update: Test POST request to /update/gauge/test1/120",
			path:   "/update/gauge/test1/120",
			method: "POST",
			body:   []byte(`{"value": 120}`),
		},
		{
			name:   "Value: Test POST request to /value/gauge/test1",
			path:   "/value/gauge/test1",
			method: "POST",
			body:   nil, // Тело запроса пустое
		},
	}

	mockStorage := storage.NewMemStorage()
	mockLogger := zap.NewExample()
	mockConfig := &config.ServerConfig{}

	fieldData := testFields{
		Storage: mockStorage,
		Logger:  mockLogger,
		Config:  mockConfig,
		// Добавьте вашу базу данных, если необходимо
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(tt.body))
			rec := httptest.NewRecorder()

			mc := &app{
				Storage: fieldData.Storage,
				Logger:  fieldData.Logger,
				Config:  fieldData.Config,
				// Добавьте вашу базу данных, если необходимо
			}

			mc.updateHandlerJSON(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("Expected status OK, got: %d", rec.Code)
			}
		})
	}
}
