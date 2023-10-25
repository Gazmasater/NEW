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

func initFlags() (*ServerConfig, *flag.FlagSet) {
	fs := flag.NewFlagSet("ServerConfig", flag.ExitOnError)

	addr := fs.String("a", "localhost:8080", "Адрес HTTP-сервера")
	storeInterval := fs.Int("i", 300, "Интервал времени в секундах для сохранения на диск")
	fileStoragePath := fs.String("f", "/tmp/metrics-db.json", "Путь к файлу для сохранения текущих значений")
	restore := fs.Bool("r", true, "Восстановление ранее сохраненных значений")
	databaseDSN := fs.String("d", "postgres://postgres:qwert@localhost:5432/postgres?sslmode=disable", "Database DSN")

	return &ServerConfig{
		Address:         *addr,
		StoreInterval:   *storeInterval,
		FileStoragePath: *fileStoragePath,
		Restore:         *restore,
		DatabaseDSN:     *databaseDSN,
	}, fs
}

func applyEnvOverrides(cfg *ServerConfig) {
	if addrEnv := os.Getenv("ADDRESS"); addrEnv != "" {
		cfg.Address = addrEnv
	}
	if storeIntervalEnv := os.Getenv("STORE_INTERVAL"); storeIntervalEnv != "" {
		if interval, err := strconv.Atoi(storeIntervalEnv); err == nil {
			cfg.StoreInterval = interval
		}
	}
	if fileStoragePathEnv := os.Getenv("FILE_STORAGE_PATH"); fileStoragePathEnv != "" {
		cfg.FileStoragePath = fileStoragePathEnv
	}
	if restoreEnv := os.Getenv("RESTORE"); restoreEnv != "" {
		if restore, err := strconv.ParseBool(restoreEnv); err == nil {
			cfg.Restore = restore
		}
	}
	if databaseDSNEnv := os.Getenv("DATABASE_DSN"); databaseDSNEnv != "" {
		cfg.DatabaseDSN = databaseDSNEnv
	}
}

func InitServerConfig() *ServerConfig {
	cfg, fs := initFlags()
	fs.Parse(os.Args[1:]) // Парсим аргументы командной строки

	applyEnvOverrides(cfg) // Применяем переменные окружения

	return cfg
}
