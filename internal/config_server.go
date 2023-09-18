package internal

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
)

// ServerConfig - структура для хранения параметров конфигурации сервера.
type ServerConfig struct {
	Address         string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
}

// InitServerConfig - функция для инициализации конфигурации server
func InitServerConfig() *ServerConfig {
	var (
		addr            string
		storeInterval   int
		fileStoragePath string
		restore         bool
	)

	addrEnv := os.Getenv("ADDRESS")
	if addrEnv != "" {
		addr = addrEnv
	} else {
		flag.StringVar(&addr, "a", "localhost:8080", "Адрес HTTP-сервера")

		if _, err := url.Parse(addr); err != nil {
			fmt.Printf("Ошибка: неверный формат адреса сервера: %s\n", addr)
			flag.Usage()
			os.Exit(1)
		}
	}

	storeIntervalEnv := os.Getenv("STORE_INTERVAL")
	if storeIntervalEnv != "" {
		storeInterval, _ = strconv.Atoi(storeIntervalEnv)
	} else {
		flag.IntVar(&storeInterval, "i", 300, "Интервал времени в секундах для сохранения на диск")
	}

	fileStoragePathEnv := os.Getenv("FILE_STORAGE_PATH")
	if fileStoragePathEnv != "" {
		fileStoragePath = fileStoragePathEnv
	} else {
		flag.StringVar(&fileStoragePath, "f", "/tmp/metrics-db.json", "Путь к файлу для сохранения текущих значений")
	}

	restoreEnv := os.Getenv("RESTORE")
	if restoreEnv != "" {
		restore, _ = strconv.ParseBool(restoreEnv)
	} else {
		flag.BoolVar(&restore, "r", true, "Восстановление ранее сохраненных значений")
	}

	flag.Parse()

	return &ServerConfig{
		Address:         addr,
		StoreInterval:   storeInterval,
		FileStoragePath: fileStoragePath,
		Restore:         restore,
	}
}
