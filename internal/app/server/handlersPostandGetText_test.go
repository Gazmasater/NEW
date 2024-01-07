package server

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"project.com/internal/storage"
)

func TestHandlePostandGetRequestGaugeText(t *testing.T) {
	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)

	// Генерируем случайное число для запроса и сохраняем его
	randomValue := rand.Float64() * 1000
	randomMetrisstring := generateRandomString()
	randomValueString := fmt.Sprintf("%.2f", randomValue)
	// Тестовые данные для обработки запросов
	tests := []struct {
		url         string
		reqBody     []byte
		contentType string
		expected    int     // Ожидаемый статус-код
		expectedGet float64 // Ожидаемое содержимое GET запроса (число с плавающей точкой)
	}{
		{
			url:         fmt.Sprintf("/update/gauge/%s/%s", randomMetrisstring, randomValueString),
			reqBody:     []byte(fmt.Sprintf("metricType=gauge&metricName=%s&metricValue=%s", randomMetrisstring, randomValueString)),
			contentType: "text/plain",
			expected:    http.StatusOK,
			expectedGet: randomValue,
		},
		// Добавьте другие тестовые сценарии по необходимости
	}

	// Создаем новый маршрутизатор Chi
	r := chi.NewRouter()

	// Создаем экземпляр обработчика
	mc := &app{
		Storage: storage.NewMemStorage(),
	}

	// Привязываем обработчик к маршруту
	r.Post("/update/{metricType}/{metricName}/{metricValue}", mc.HandlePostRequest)
	r.Get("/value/{metricType}/{metricName}", mc.HandleGetRequest)

	// Создаем тестовый сервер
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(fmt.Sprintf("URL: %s", tt.url), func(t *testing.T) {
			// Создаем тестовый сервер с использованием маршрутизатора Chi

			// Формируем полный URL для тестирования POST-запроса
			url := ts.URL + tt.url
			// Создаем тестовый запрос
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(tt.reqBody))
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
			getURL := ts.URL + fmt.Sprintf("/value/gauge/%s", randomMetrisstring)
			getReq, err := http.NewRequest("GET", getURL, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Выполняем GET запрос к тестовому серверу
			getRes, err := http.DefaultClient.Do(getReq)
			if err != nil {
				t.Fatal(err)
			}
			defer getRes.Body.Close()

			// Проверяем статус-код GET ответа
			if getRes.StatusCode != http.StatusOK {
				t.Errorf("Expected status %d; got %d", http.StatusOK, getRes.StatusCode)
			}

			// Считываем и преобразуем ответ GET запроса в число с плавающей точкой
			var getBodyFloat float64
			if _, err := fmt.Fscanf(getRes.Body, "%f", &getBodyFloat); err != nil {
				t.Fatal(err)
			}

			// Указываем относительную погрешность в пределах 0.01 (0.1%)
			epsilon := 0.001 * tt.expectedGet

			// Проверяем, что разница между значениями не превышает относительную погрешность
			if math.Abs(tt.expectedGet-getBodyFloat) > epsilon {
				t.Errorf("Expected body %.2f; got %.2f", tt.expectedGet, getBodyFloat)
			}
		})
	}
}

func TestHandlePostandGetRequestCounterText(t *testing.T) {
	source := rand.NewSource(time.Now().UnixNano())

	rand := rand.New(source)
	var sum int

	// Генерируем случайное число для запроса и сохраняем его
	randomValue1 := rand.Intn(1001)
	randomValueString1 := strconv.Itoa(randomValue1)
	randomMetricstring1 := generateRandomString()

	randomValue2 := rand.Intn(1001)
	randomValueString2 := strconv.Itoa(randomValue2)

	randomValue3 := rand.Intn(1001)
	randomValueString3 := strconv.Itoa(randomValue3)

	// Тестовые данные для обработки запросов
	tests := []struct {
		url         string
		reqBody     []byte
		contentType string
		expected    int    // Ожидаемый статус-код
		expectedGet string // Ожидаемое содержимое GET запроса
		metricName  string
	}{
		{
			url:         fmt.Sprintf("/update/counter/%s/%d", randomMetricstring1, randomValue1),
			reqBody:     []byte(fmt.Sprintf("metricType=counter&metricName=%s&metricValue=%s", randomMetricstring1, randomValueString1)),
			contentType: "text/plain",
			expected:    http.StatusOK,
			expectedGet: randomValueString1,
			metricName:  randomMetricstring1,
		},

		{
			url:         fmt.Sprintf("/update/counter/%s/%d", randomMetricstring1, randomValue2),
			reqBody:     []byte(fmt.Sprintf("metricType=counter&metricName=%s&metricValue=%s", randomMetricstring1, randomValueString2)),
			contentType: "text/plain",
			expected:    http.StatusOK,
			expectedGet: randomValueString2,
			metricName:  randomMetricstring1,
		},

		{
			url:         fmt.Sprintf("/update/counter/%s/%d", randomMetricstring1, randomValue3),
			reqBody:     []byte(fmt.Sprintf("metricType=counter&metricName=%s&metricValue=%s", randomMetricstring1, randomValueString3)),
			contentType: "text/plain",
			expected:    http.StatusOK,
			expectedGet: randomValueString3,
			metricName:  randomMetricstring1,
		},
		// Добавьте другие тестовые сценарии по необходимости
	}

	// Создаем новый маршрутизатор Chi
	r := chi.NewRouter()

	// Создаем экземпляр обработчика
	mc := &app{
		Storage: storage.NewMemStorage(),
	}

	// Привязываем обработчик к маршруту
	r.Post("/update/{metricType}/{metricName}/{metricValue}", mc.HandlePostRequest)
	r.Get("/value/{metricType}/{metricName}", mc.HandleGetRequest)

	// Создаем тестовый сервер
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(fmt.Sprintf("URL: %s", tt.url), func(t *testing.T) {

			// Формируем полный URL для тестирования POST-запроса
			url := ts.URL + tt.url

			// Создаем тестовый запрос
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(tt.reqBody))
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

			// Парсим значение из ответа на пост-запрос
			value, err := strconv.Atoi(tt.expectedGet)
			if err != nil {
				t.Fatal(err)
			}

			// Добавляем значение к сумме
			sum += value

			// Создаем GET запрос
			getURL := ts.URL + "/value/counter/" + tt.metricName

			getReq, err := http.NewRequest("GET", getURL, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Выполняем GET запрос к тестовому серверу
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

			sumString := strconv.Itoa(sum)

			if getBody != sumString {
				t.Errorf("Expected body %s; got %s", sumString, getBody)
			}
		})
	}
}
