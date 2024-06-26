package config

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Config struct {
	Address        string `json:"address"`
	Port           uint16 `json:"port"`
	Addition       uint
	Subtraction    uint
	Multiplication uint
	Division       uint
}

// GetAddress возвращает полный адрес
func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

// CustomFormatter определяет свой собственный формат вывода для логгера
type CustomFormatter struct{}

// Format форматирует запись лога с заданным форматом времени
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s] [%s] %s\n",
		entry.Time.Format("15:04:05.0000000"), // Формат времени: часы:минуты:секунды.микросекунды
		strings.ToUpper(entry.Level.String()),
		entry.Message,
	)), nil
}

// LoadConfig принимает порт для сервера и длительность математических операций и возвращает конфиг
func LoadConfig(port uint16, addition, subtraction, multiplication, division uint) (*Config, error) {
	if port <= 0 || port > 65535 {
		return nil, errors.New("invalid port")
	}

	return &Config{
		Address:        "localhost",
		Port:           port,
		Addition:       addition,
		Subtraction:    subtraction,
		Multiplication: multiplication,
		Division:       division,
	}, nil
}

func LoadLogger(debugLevel bool) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&CustomFormatter{})
	log.SetOutput(os.Stdout)
	if debugLevel {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}
	return log
}
