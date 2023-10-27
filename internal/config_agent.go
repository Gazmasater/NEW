package internal

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
)

// AgentConfig - структура для хранения параметров конфигурации агента.
type AgentConfig struct {
	Address        string
	ReportInterval int
	PollInterval   int
	Key            string
}

// InitAgentConfig - функция для инициализации конфигурации агента.
func InitAgentConfig() *AgentConfig {
	println("InitAgentConfig internal")
	var (
		reportSeconds int
		pollSeconds   int
		addr          string
		key           string
	)

	addrEnv := os.Getenv("ADDRESS")
	if addrEnv != "" {
		addr = addrEnv
	}
	flag.StringVar(&addr, "a", "localhost:8080", "Адрес HTTP-сервера")
	if _, err := url.Parse(addr); err != nil {
		fmt.Printf("Ошибка: неверный формат адреса сервера в модуле агента: %s\n", addr)
		return nil
	}

	reportSecondsEnv := os.Getenv("REPORT_INTERVAL")
	if reportSecondsEnv != "" {
		reportSeconds, _ = strconv.Atoi(reportSecondsEnv)
	} else {
		flag.IntVar(&reportSeconds, "r", 10, "Частота отправки метрик на сервер (в секундах)")
		if reportSeconds <= 0 {
			fmt.Println("Частота отправки метрик должна быть положительным числом.")
			flag.Usage()
			return nil
		}
	}

	pollSecondsEnv := os.Getenv("POLL_INTERVAL")
	if pollSecondsEnv != "" {
		pollSeconds, _ = strconv.Atoi(pollSecondsEnv)
	} else {
		flag.IntVar(&pollSeconds, "p", 2, "Частота опроса метрик из пакета runtime (в секундах)")
		if pollSeconds <= 0 {
			fmt.Println("Частота опроса метрик должна быть положительным числом.")
			flag.Usage()
			return nil
		}
	}

	keyEnv := os.Getenv("KEY")
	if keyEnv != "" {
		key = keyEnv
	} else {
		flag.StringVar(&key, "k", "MyKey", "Ключ для подписи данных")
	}

	flag.Parse()

	return &AgentConfig{
		Address:        addr,
		ReportInterval: reportSeconds,
		PollInterval:   pollSeconds,
		Key:            key,
	}
}
