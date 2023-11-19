package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"go.uber.org/zap"
	"project.com/internal/logger"
)

// AgentConfig - структура для хранения параметров конфигурации агента.
type AgentConfig struct {
	Address        string
	ReportInterval int
	PollInterval   int
	Key            string
	RateLimit      int // Изменил тип поля на int
}

var cfga *AgentConfig

// New - функция для инициализации конфигурации агента.
func New() (*AgentConfig, error) {
	var (
		reportSeconds int
		pollSeconds   int
		addr          string
		key           string
		rateLimit     int
	)
	log, _ := logger.New()

	flag.StringVar(&addr, "a", "localhost:8080", "Адрес HTTP-сервера")
	if _, err := url.Parse(addr); err != nil {
		log.Error("Ошибка при обработке данных localhost", zap.Error(err))

		return nil, err
	}

	addrEnv := os.Getenv("ADDRESS")
	if addrEnv != "" {
		addr = addrEnv
	}

	flag.IntVar(&reportSeconds, "r", 10, "Частота отправки метрик на сервер (в секундах)")
	if reportSeconds <= 0 {
		flag.Usage()
		log.Fatal("частота отправки метрик должна быть положительным числом")
	}

	reportSecondsEnv := os.Getenv("REPORT_INTERVAL")
	if reportSecondsEnv != "" {
		reportSeconds, _ = strconv.Atoi(reportSecondsEnv)
	}

	flag.IntVar(&pollSeconds, "p", 2, "Частота опроса метрик из пакета runtime (в секундах)")
	if pollSeconds <= 0 {
		flag.Usage()
		log.Fatal("частота опроса метрик должна быть положительным числом")
	}

	pollSecondsEnv := os.Getenv("POLL_INTERVAL")
	if pollSecondsEnv != "" {
		pollSeconds, _ = strconv.Atoi(pollSecondsEnv)
	}

	flag.StringVar(&key, "k", "MyKey", "Ключ для подписи данных")

	keyEnv := os.Getenv("KEY")
	if keyEnv != "" {
		key = keyEnv
	}

	flag.IntVar(&rateLimit, "l", 0, "Rate Limit")
	flag.IntVar(&rateLimit, "RATE_LIMIT", 0, "Rate Limit (переменная окружения)")

	flag.Parse()

	rateLimitStr := strconv.Itoa(rateLimit)
	rateLimitInt, err := strconv.Atoi(rateLimitStr)
	if err != nil {

		fmt.Println("Ошибка при преобразовании строки в число:", err)

	}

	cfga = &AgentConfig{
		Address:        addr,
		ReportInterval: reportSeconds,
		PollInterval:   pollSeconds,
		Key:            key,
		RateLimit:      rateLimitInt,
	}
	return cfga, nil
}
