package internal

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

// ServerConfig - структура для хранения параметров конфигурации сервера.
type ServerConfig struct {
	Address string
}

// InitServerConfig - функция для инициализации конфигурации сервера.
func InitServerConfig() *ServerConfig {
	var addr string

	// Чтение переменной окружения или установка значения по умолчанию
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

	flag.Parse()

	return &ServerConfig{
		Address: addr,
	}
}