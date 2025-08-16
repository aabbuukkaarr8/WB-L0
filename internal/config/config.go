package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LogLevel string       `yaml:"loglevel"`
	Server   ServerConfig `yaml:"server"`
	DB       DBConfig     `yaml:"db"`
	Kafka    KafkaConfig  `yaml:"kafka"`
}
type ServerConfig struct {
	Port string `yaml:"port"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Login, c.Password, c.Host, c.Port, c.DBName, c.SSLMode,
	)
}

type KafkaConfig struct {
	Broker  string `yaml:"broker"`
	Topic   string `yaml:"topic"`
	GroupID string `yaml:"group_id"`
}

func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
