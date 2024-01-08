package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"reflect"
	"strconv"

	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"

	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi"
	"project.com/internal/config"
	"project.com/internal/models"
	"project.com/internal/storage"
)

func TestHandlePostandGetRequestCounterJSon(t *testing.T) {

	var metrP models.Metrics
	var metrG models.Metrics
	var sum int64

	serverCfg := config.InitServerConfig()

	db, err := sql.Open("postgres", serverCfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Ошибка при открытии соединения с базой данных: %v", err)
		return
	}
	defer db.Close()

	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)

	// Генерируем случайное число для запроса и сохраняем его
	randomValue1 := rand.Intn(1001)
	randomValueString1 := strconv.Itoa(randomValue1)

	randomValue2 := rand.Intn(1001)
	randomValueString2 := strconv.Itoa(randomValue2)

	randomValue3 := rand.Intn(1001)
	randomValueString3 := strconv.Itoa(randomValue3)

	randomMetricstring := generateRandomString()

	// Тестовые данные для обработки запросов
	tests := []struct {
		url         string
		reqBodyPost []byte
		reqBodyGet  []byte
		contentType string
		expected    int    // Ожидаемый статус-код
		expectedGet string // Ожидаемое содержимое GET запроса
	}{
		{
			url:         "/update/",
			reqBodyPost: []byte(fmt.Sprintf(`{"Type":"counter","ID":"%s","Delta":%s}`, randomMetricstring, randomValueString1)),
			reqBodyGet:  []byte(fmt.Sprintf(`{"Type":"counter","ID":"%s"}`, randomMetricstring)),
			contentType: "application/json",
			expected:    http.StatusOK,
			expectedGet: fmt.Sprintf(`{"type":"counter","id":"%s","delta":%s}`, randomMetricstring, randomValueString1),
		},

		{
			url:         "/update/",
			reqBodyPost: []byte(fmt.Sprintf(`{"Type":"counter","ID":"%s","Delta":%s}`, randomMetricstring, randomValueString2)),
			reqBodyGet:  []byte(fmt.Sprintf(`{"Type":"counter","ID":"%s"}`, randomMetricstring)),
			contentType: "application/json",
			expected:    http.StatusOK,
			expectedGet: fmt.Sprintf(`{"type":"counter","id":"%s","delta":%s}`, randomMetricstring, randomValueString2),
		},

		{
			url:         "/update/",
			reqBodyPost: []byte(fmt.Sprintf(`{"Type":"counter","ID":"%s","Delta":%s}`, randomMetricstring, randomValueString3)),
			reqBodyGet:  []byte(fmt.Sprintf(`{"Type":"counter","ID":"%s"}`, randomMetricstring)),
			contentType: "application/json",
			expected:    http.StatusOK,
			expectedGet: fmt.Sprintf(`{"type":"counter","id":"%s","delta":%s}`, randomMetricstring, randomValueString3),
		},
		// Добавьте другие тестовые сценарии по необходимости
	}

	// Создаем новый маршрутизатор Chi
	r := chi.NewRouter()

	// Создаем экземпляр обработчика
	mc := &app{
		Storage: storage.NewMemStorage(),
		Config:  serverCfg,
		DB:      db,
	}

	// Привязываем обработчик к маршруту
	r.Post("/update/", mc.updateHandlerJSONOptimiz)
	r.Post("/value/", mc.updateHandlerJSONValueOptimiz)

	// Создаем тестовый сервер
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(fmt.Sprintf("URL: %s", tt.url), func(t *testing.T) {

			// Формируем полный URL для тPOST-запроса
			url := ts.URL + tt.url

			// Создаем тестовый запрос
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(tt.reqBodyPost))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", tt.contentType)

			// Выполняем запрос к тестовому серверу
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			// Проверяем статус-код ответа
			if res.StatusCode != tt.expected {
				t.Errorf("Expected status %d; got %d", tt.expected, res.StatusCode)
			}

			// Создаем GET запрос
			getURL := ts.URL + "/value/"

			getReq, err := http.NewRequest("POST", getURL, bytes.NewBuffer(tt.reqBodyGet))
			if err != nil {
				t.Fatal(err)
			}

			// Устанавливаем заголовок Content-Type для JSON
			getReq.Header.Set("Content-Type", tt.contentType)

			// Выполняем GET(он же POST) запрос к тестовому серверу
			getRes, err := http.DefaultClient.Do(getReq)
			if err != nil {
				t.Fatal(err)
			}
			defer getRes.Body.Close()

			body, err := io.ReadAll(getRes.Body)
			if err != nil {
				t.Fatal(err)
			}

			// Проверяем статус-код GET ответа
			if getRes.StatusCode != http.StatusOK {
				t.Errorf("Expected status %d; got %d", http.StatusOK, getRes.StatusCode)
			}

			if err := json.Unmarshal(tt.reqBodyPost, &metrP); err != nil {
				fmt.Println("Ошибка при парсинге JSON:", err)
				return
			}

			sum += *metrP.Delta

			if err := json.Unmarshal(body, &metrG); err != nil {
				fmt.Println("Ошибка при парсинге JSON:", err)
				return
			}

			if sum != *metrG.Delta {
				t.Errorf("Expected body %s; got %d", tt.reqBodyPost, *metrG.Delta)

			}

		})
	}
}

func TestHandlePostandGetRequestGaugeJson(t *testing.T) {

	serverCfg := config.InitServerConfig()

	println("serverCfg.DatabaseDSN", serverCfg.DatabaseDSN)

	db, err := sql.Open("postgres", serverCfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Ошибка при открытии соединения с базой данных: %v", err)
		return
	}
	defer db.Close()

	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)

	// Генерируем случайное число для запроса и сохраняем его
	randomValue := rand.Float64() * 1000
	randomValueString := fmt.Sprintf("%.2f", randomValue)
	randomMetricstring := generateRandomString()

	// Тестовые данные для обработки запросов
	tests := []struct {
		url         string
		reqBodyPost []byte
		reqBodyGet  []byte
		contentType string
		expected    int    // Ожидаемый статус-код
		expectedGet string // Ожидаемое содержимое GET запроса
	}{
		{
			url:         "/update/",
			reqBodyPost: []byte(fmt.Sprintf(`{"Type":"gauge","ID":"%s","Value":%s}`, randomMetricstring, randomValueString)),
			reqBodyGet:  []byte(fmt.Sprintf(`{"Type":"gauge","ID":"%s"}`, randomMetricstring)),
			contentType: "application/json",
			expected:    http.StatusOK,
			expectedGet: fmt.Sprintf(`{"type":"gauge","id":"%s","value":%s}`, randomMetricstring, randomValueString),
		},
		// Добавьте другие тестовые сценарии по необходимости
	}

	// Создаем новый маршрутизатор Chi
	r := chi.NewRouter()

	// Создаем экземпляр обработчика
	mc := &app{
		Storage: storage.NewMemStorage(),
		Config:  serverCfg,
		DB:      db,
	}

	// mc.SetupDatabase()

	// Привязываем обработчик к маршруту
	r.Post("/update/", mc.updateHandlerJSONOptimiz)
	r.Post("/value/", mc.updateHandlerJSONValueOptimiz)

	// Создаем тестовый сервер
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(fmt.Sprintf("URL: %s", tt.url), func(t *testing.T) {

			// Формируем полный URL для тPOST-запроса
			url := ts.URL + tt.url

			fmt.Println("POST Request Body:", string(tt.reqBodyPost))
			// Создаем тестовый запрос
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(tt.reqBodyPost))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", tt.contentType)

			// Выполняем запрос к тестовому серверу
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			// Чтение тела ответа

			// Проверяем статус-код ответа
			if res.StatusCode != tt.expected {
				t.Errorf("Expected status %d; got %d", tt.expected, res.StatusCode)
			}

			// Создаем GET запрос
			getURL := ts.URL + "/value/"

			getReq, err := http.NewRequest("POST", getURL, bytes.NewBuffer(tt.reqBodyGet))
			if err != nil {
				t.Fatal(err)
			}

			// Устанавливаем заголовок Content-Type для JSON
			getReq.Header.Set("Content-Type", tt.contentType)

			// Выполняем POST запрос к тестовому серверу
			getRes, err := http.DefaultClient.Do(getReq)
			if err != nil {
				t.Fatal(err)
			}
			defer getRes.Body.Close()

			// Проверяем статус-код GET ответа
			if getRes.StatusCode != http.StatusOK {
				t.Errorf("Expected status %d; got %d", http.StatusOK, getRes.StatusCode)
			}

			// Сравниваем тело GET ответа с ожидаемым результатом
			buf := new(bytes.Buffer)
			buf.ReadFrom(getRes.Body)
			getBody := buf.String()
			fmt.Println("GET Response Body Gauge:", getBody)
			var expectedMap map[string]interface{}
			if err := json.Unmarshal([]byte(tt.expectedGet), &expectedMap); err != nil {
				t.Fatalf("Failed to unmarshal expected JSON: %s", err)
			}

			var bodyMap map[string]interface{}
			if err := json.Unmarshal([]byte(getBody), &bodyMap); err != nil {
				t.Fatalf("Failed to unmarshal body JSON: %s", err)
			}

			if !reflect.DeepEqual(expectedMap, bodyMap) {
				t.Errorf("Expected JSON %v; got %v", expectedMap, bodyMap)
			}
		})
	}
}
