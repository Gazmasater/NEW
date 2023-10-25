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

	addrEnv := os.Getenv("ADDRESS")
	if addrEnv == "" {
		addr = "localhost:8080"
	}
	flag.StringVar(&addr, "a", addrEnv, "Адрес HTTP-сервера")

	fileStoragePathEnv := os.Getenv("FILE_STORAGE_PATH")
	if fileStoragePathEnv == "" {
		fileStoragePath = "/tmp/metrics-db.json"
	}
	flag.StringVar(&fileStoragePath, "f", fileStoragePathEnv, "Путь к файлу для сохранения текущих значений")

	databaseDSNEnv := os.Getenv("DATABASE_DSN")
	if databaseDSNEnv == "" {
		databaseDSN = "postgres://postgres:qwert@localhost:5432/postgres?sslmode=disable"
	}
	flag.StringVar(&databaseDSN, "d", databaseDSNEnv, "Database DSN")

	flag.BoolVar(&restore, "r", true, "Восстановление ранее сохраненных значений")
	flag.IntVar(&storeInterval, "i", 300, "Интервал времени в секундах для сохранения на диск")

	storeIntervalEnv := os.Getenv("STORE_INTERVAL")
	if storeIntervalEnv != "" {
		storeInterval, _ = strconv.Atoi(storeIntervalEnv)
	}

	restoreEnv := os.Getenv("RESTORE")
	if restoreEnv != "" {
		restore, _ = strconv.ParseBool(restoreEnv)
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
