package server

import (
	"bytes"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"project.com/internal/config"
)

func TestApp_updateHandlerJSON(t *testing.T) {

	db, err := sql.Open("postgres", mc.Config.DatabaseDSN)
	if err != nil {
		log.Printf("Ошибка при открытии базы данных: %v", err)
		return err
	}
	defer db.Close()
	// Подготовка фейковых данных для тестов
	mockStorage := storage.NewMemStorage() // Используйте вашу реализацию хранилища
	mockLogger := zap.NewExample()
	mockConfig := &config.ServerConfig{} // Подставьте реальную конфигурацию

	// Создание фейкового HTTP запроса
	req := httptest.NewRequest("POST", "/update/", bytes.NewBufferString(`{"id": "test1", "type": "counter", "delta": 10}`))
	rr := httptest.NewRecorder()

	// Создание экземпляра вашего приложения
	mc := &app{
		Storage: mockStorage,
		Logger:  mockLogger,
		Config:  mockConfig,
		DB:      db, // Добавьте вашу базу данных, если используете
	}

	// Вызов функции обработчика
	mc.updateHandlerJSON(rr, req)

	// TODO: Добавьте утверждения, чтобы проверить ожидаемое поведение в ответ на запрос.
	// Проверьте, что метрика была успешно обработана или что была возвращена ошибка в случае неудачи.
	// Например:
	// - Проверка статус кода ответа (rr.Code)
	// - Проверка содержимого ответа (rr.Body)
	// - Проверка изменений в вашем хранилище или файлах

	// Пример утверждения:
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got: %d", rr.Code)
	}

	// ... Добавьте другие утверждения по вашему усмотрению
}
