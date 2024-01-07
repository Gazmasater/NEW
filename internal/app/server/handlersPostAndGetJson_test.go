package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi"
	"project.com/internal/config"
	"project.com/internal/storage"
)

// func TestHandlePostandGetRequestCounterJSon(t *testing.T) {
// 	source := rand.NewSource(time.Now().UnixNano())

// 	// Инициализируем генератор случайных чисел с использованием источника
// 	rand := rand.New(source)
// 	var sum int

// 	// Генерируем случайное число для запроса и сохраняем его
// 	randomValue1 := rand.Intn(1001)
// 	randomValueString1 := strconv.Itoa(randomValue1)
// 	randomValue2 := rand.Intn(1001)
// 	randomValueString2 := strconv.Itoa(randomValue2)
// 	randomValue3 := rand.Intn(1001)
// 	randomValueString3 := strconv.Itoa(randomValue3)

// 	// Тестовые данные для обработки запросов
// 	tests := []struct {
// 		url         string
// 		reqBody     []byte
// 		contentType string
// 		expected    int    // Ожидаемый статус-код
// 		expectedGet string // Ожидаемое содержимое GET запроса
// 	}{
// 		{
// 			url:         fmt.Sprintf("/update/counter/test1/%s", randomValueString1),
// 			reqBody:     []byte(fmt.Sprintf(`{"MType":"counter","ID":"test1","Delta":%s}`, randomValueString1)),
// 			contentType: "application/json",
// 			expected:    http.StatusOK,
// 			expectedGet: randomValueString1,
// 		},

// 		{
// 			url:         fmt.Sprintf("/update/counter/test1/%s", randomValueString2),
// 			reqBody:     []byte(fmt.Sprintf(`{"MType":"counter","ID":"test1","Delta":%s}`, randomValueString2)),
// 			contentType: "application/json",
// 			expected:    http.StatusOK,
// 			expectedGet: randomValueString2,
// 		},

// 		{
// 			url:         fmt.Sprintf("/update/counter/test1/%s", randomValueString3),
// 			reqBody:     []byte(fmt.Sprintf(`{"MType":"counter","ID":"test1","Delta":%s}`, randomValueString3)),
// 			contentType: "application/json",
// 			expected:    http.StatusOK,
// 			expectedGet: randomValueString3,
// 		},
// 		// Добавьте другие тестовые сценарии по необходимости
// 	}

// 	// Создаем новый маршрутизатор Chi
// 	r := chi.NewRouter()

// 	// Создаем экземпляр обработчика
// 	mc := &app{
// 		Storage: storage.NewMemStorage(),
// 	}

// 	// Привязываем обработчик к маршруту
// 	r.Post("/update/", mc.updateHandlerJSON)
// 	r.Post("/value/", mc.updateHandlerJSONValue)

// 	for _, tt := range tests {
// 		t.Run(fmt.Sprintf("URL: %s", tt.url), func(t *testing.T) {
// 			// Создаем тестовый сервер с использованием маршрутизатора Chi

// 			// Формируем полный URL для тестирования POST-запроса
// 			url := "http://localhost:8080" + tt.url

// 			// Создаем тестовый запрос
// 			req, err := http.NewRequest("POST", url, bytes.NewBuffer(tt.reqBody))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Set("Content-Type", tt.contentType)

// 			fmt.Println("POST Request Body Json Counter:", string(tt.reqBody))

// 			// Выполняем запрос к тестовому серверу
// 			res, err := http.DefaultClient.Do(req)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			defer res.Body.Close()

// 			// Проверяем статус-код ответа
// 			if res.StatusCode != tt.expected {
// 				t.Errorf("Expected status %d; got %d", tt.expected, res.StatusCode)
// 			}

// 			// Парсим значение из ответа на пост-запрос
// 			value, err := strconv.Atoi(tt.expectedGet)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			// Добавляем значение к сумме
// 			sum += value

// 			// Создаем GET запрос
// 			getURL := "http://localhost:8080" + "/value/counter/test1"
// 			getReq, err := http.NewRequest("GET", getURL, nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			// Выполняем GET запрос к тестовому серверу
// 			getRes, err := http.DefaultClient.Do(getReq)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			defer getRes.Body.Close()

// 			// Проверяем статус-код GET ответа
// 			if getRes.StatusCode != http.StatusOK {
// 				t.Errorf("Expected status %d; got %d", http.StatusOK, getRes.StatusCode)
// 			}

// 			// Сравниваем тело GET ответа с ожидаемым результатом
// 			buf := new(bytes.Buffer)
// 			buf.ReadFrom(getRes.Body)
// 			getBody := buf.String()
// 			fmt.Println("GET Response Body Json counter:", getBody)

// 			sumString := strconv.Itoa(sum)
// 			if getBody != sumString {
// 				t.Errorf("Expected body %s; got %s", sumString, getBody)
// 			}
// 		})
// 	}
// }

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
			reqBodyPost: []byte(fmt.Sprintf(`{"Type":"gauge","ID":"test1","Value":%s}`, randomValueString)),
			reqBodyGet:  []byte(`{"Type":"gauge","ID":"test1"}`),
			contentType: "application/json",
			expected:    http.StatusOK,
			expectedGet: fmt.Sprintf(`{"id":"test1","type":"gauge","value":%s}`, randomValueString),
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
	r.Post("/update/", mc.updateHandlerJSON)
	r.Post("/value/", mc.updateHandlerJSONValue)

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
