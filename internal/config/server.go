package config

import (
	"flag"
	"os"
	"strconv"
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

	flag.StringVar(&addr, "a", "localhost:8080", "Адрес HTTP-сервера")
	flag.IntVar(&storeInterval, "i", 300, "Интервал времени в секундах для сохранения на диск")
	flag.StringVar(&fileStoragePath, "f", "/tmp/metrics-db.json", "Путь к файлу для сохранения текущих значений")
	flag.BoolVar(&restore, "r", true, "Восстановление ранее сохраненных значений")
	flag.StringVar(&databaseDSN, "dbdsn", "postgres://postgres:qwert@localhost:5432/postgres?sslmode=disable", "Database DSN")

	// Проверяем переменные окружения и используем их, если они определены
	if addrEnv := os.Getenv("ADDRESS"); addrEnv != "" {
		addr = addrEnv
	}

	if storeIntervalEnv := os.Getenv("STORE_INTERVAL"); storeIntervalEnv != "" {
		storeInterval, _ = strconv.Atoi(storeIntervalEnv)
	}

	if fileStoragePathEnv := os.Getenv("FILE_STORAGE_PATH"); fileStoragePathEnv != "" {
		fileStoragePath = fileStoragePathEnv
	}

	if restoreEnv := os.Getenv("RESTORE"); restoreEnv != "" {
		restore, _ = strconv.ParseBool(restoreEnv)
	}

	if databaseDSNEnv := os.Getenv("DATABASE_DSN"); databaseDSNEnv != "" {
		databaseDSN = databaseDSNEnv
	}

	flag.Parse()

	return &ServerConfig{
		Address:         addr,
		StoreInterval:   storeInterval,
		FileStoragePath: fileStoragePath,
		Restore:         restore,
		DatabaseDSN:     databaseDSN,
	}
}
