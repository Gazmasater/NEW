package server

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
	"project.com/internal/models"
)

func (mc *app) SetupDatabase() error {
	// Открываем соединение с базой данных
	db, err := sql.Open("postgres", mc.Config.DatabaseDSN)
	if err != nil {
		log.Printf("Ошибка при открытии базы данных: %v", err)
		return err
	}
	defer db.Close()

	// Запрос для создания таблицы
	createTableQuery := `
        CREATE TABLE IF NOT EXISTS metrics (
            name VARCHAR(255) NOT NULL,
            type VARCHAR(50) NOT NULL,
            value DOUBLE PRECISION,
            delta BIGINT
        )
    `

	// Выполняем запрос для создания таблицы
	_, err = db.Exec(createTableQuery)
	if err != nil {
		pqErr, isPQError := err.(*pq.Error)
		if isPQError && pqErr.Code == "23505" {
			// Код "23505" соответствует ошибке уникального нарушения.
			log.Printf("Ошибка при создании таблицы: %v (UniqueViolation)", err)
			return err
		}
		log.Printf("Ошибка при создании таблицы: %v", err)
		return err
	}

	return nil
}

func (mc *app) WriteMetricToDatabase(metric models.Metrics) error {
	var query string
	var args []any

	switch metric.MType {
	case "gauge":
		query = "INSERT INTO metrics (name, type, value) VALUES ($1, $2, $3)"
		args = []interface{}{metric.ID, metric.MType, metric.Value}
	case "counter":
		query = "INSERT INTO metrics (name, type, delta) VALUES ($1, $2, $3)"
		args = []interface{}{metric.ID, metric.MType, metric.Delta}
	default:
		log.Printf("Неизвестный тип метрики: %s", metric.MType)
		return fmt.Errorf("неизвестный тип метрики")

	}
	if mc.DB == nil {
		log.Println("Ошибка: mc.DB не инициализирован.")
		return fmt.Errorf("mc.DB не инициализирован")
	}

	// Проверяем, существует ли метрика с такими же значениями name и type
	var count int
	err := mc.DB.QueryRow("SELECT COUNT(*) FROM metrics WHERE name = $1 AND type = $2", metric.ID, metric.MType).Scan(&count)
	if err != nil {
		log.Printf("Ошибка при проверке наличия метрики: %s", err)
		return err
	}

	if count > 0 {
		// Метрика с такими значениями name и type существует, удаляем ее
		_, err := mc.DB.Exec("DELETE FROM metrics WHERE name = $1 AND type = $2", metric.ID, metric.MType)
		if err != nil {
			log.Printf("Ошибка при удалении метрики: %s", err)
			return err
		}
	}

	// Теперь выполняем вставку новой метрики
	_, err = mc.DB.Exec(query, args...)
	if err != nil {
		log.Printf("Ошибка при записи метрики в базу данных: %s", err)
		return err
	}
	return nil
}
