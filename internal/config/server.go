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
		addrEnv = "localhost:8080"
	}
	flag.StringVar(&addr, "a", addrEnv, "Адрес HTTP-сервера")

	storeIntervalEnv := os.Getenv("STORE_INTERVAL")
	if storeIntervalEnv == "" {
		storeIntervalEnv = "300"
	}
	storeInterval, _ = strconv.Atoi(storeIntervalEnv) // Преобразование строки в int
	flag.IntVar(&storeInterval, "i", storeInterval, "Интервал времени в секундах для сохранения на диск")

	fileStoragePathEnv := os.Getenv("FILE_STORAGE_PATH")
	if fileStoragePathEnv == "" {
		fileStoragePath = "/tmp/metrics-db.json"
	}
	flag.StringVar(&fileStoragePath, "f", fileStoragePath, "Путь к файлу для сохранения текущих значений")

	restoreEnv := os.Getenv("RESTORE")
	if restoreEnv == "" {
		restore = true // Присваиваем булевое значение напрямую
	}
	flag.BoolVar(&restore, "r", restore, "Восстановление ранее сохраненных значений")

	databaseDSNEnv := os.Getenv("DATABASE_DSN")
	if databaseDSNEnv != "" {
		databaseDSN = "postgres://postgres:qwert@localhost:5432/postgres?sslmode=disable"
	}
	flag.StringVar(&databaseDSN, "d", databaseDSNEnv, "Database DSN")

	flag.Parse()

	return &ServerConfig{
		Address:         addr,
		StoreInterval:   storeInterval,
		FileStoragePath: fileStoragePath,
		Restore:         restore,
		DatabaseDSN:     databaseDSN,
	}
}
