package server

import (
	"database/sql"
	"log"

	"github.com/lib/pq"
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
