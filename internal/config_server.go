package internal

import (
	"flag"
)

type ServerConfig struct {
	Address         string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	DatabaseDSN     string
}

func InitServerConfig() *ServerConfig {
	var (
		addr            string
		storeInterval   int
		fileStoragePath string
		restore         bool
		databaseDSN     string
	)

	// Установка значений по умолчанию
	flag.StringVar(&addr, "a", "localhost:8080", "Адрес HTTP-сервера")
	flag.IntVar(&storeInterval, "i", 300, "Интервал времени в секундах для сохранения на диск")
	flag.StringVar(&fileStoragePath, "f", "/tmp/metrics-db.json", "Путь к файлу для сохранения текущих значений")
	flag.BoolVar(&restore, "r", true, "Восстановление ранее сохраненных значений")
	flag.StringVar(&databaseDSN, "dbdsn", "postgres://postgres:qwert@localhost:5432/postgres?sslmode=disable", "Database DSN")

	flag.Parse()

	// Оставляем только создание ServerConfig с заданными значениями
	return &ServerConfig{
		Address:         addr,
		StoreInterval:   storeInterval,
		FileStoragePath: fileStoragePath,
		Restore:         restore,
		DatabaseDSN:     databaseDSN,
	}
}
