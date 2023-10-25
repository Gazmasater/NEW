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
	flag.StringVar(&databaseDSN, "d", "postgres://postgres:qwert@localhost:5432/postgres?sslmode=disable", "Database DSN")

	// Проверяем переменные окружения и используем их, если они определены
	addrEnv := os.Getenv("ADDRESS")
	if addrEnv != "" {
		addr = addrEnv
	}

	storeIntervalEnv := os.Getenv("STORE_INTERVAL")
	if storeIntervalEnv != "" {
		storeInterval, _ = strconv.Atoi(storeIntervalEnv)
	}

	fileStoragePathEnv := os.Getenv("FILE_STORAGE_PATH")
	if fileStoragePathEnv != "" {
		fileStoragePath = fileStoragePathEnv
	}

	restoreEnv := os.Getenv("RESTORE")
	if restoreEnv != "" {
		restore, _ = strconv.ParseBool(restoreEnv)
	}

	databaseDSNEnv := os.Getenv("DATABASE_DSN")
	if databaseDSNEnv != "" {
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
