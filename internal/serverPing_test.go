package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"

	"go.uber.org/zap"
)

func TestHandlerDependencies_Ping(t *testing.T) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	config := &ServerConfig{
		Address:     "localhost:8080",
		DatabaseDSN: "postgres://postgres:qwert@localhost:5432/postgres?sslmode=disable",
	}

	mc := &HandlerDependencies{
		Storage: &MemStorage{},
		Logger:  zap.NewNop(),
		Config:  config,
	}

	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Вызываем функцию Ping
	mc.Ping(rr, req)

	// Проверка HTTP-статуса ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус %v, получено %v", http.StatusOK, status)
	}

	// Проверка тела ответа
	expectedResponse := "Database is working\n"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Ожидалось тело ответа '%s', получено '%s'", expectedResponse, rr.Body.String())
	}

}
