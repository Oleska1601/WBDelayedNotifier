package config

import (
	"fmt"
	"time"

	wbfcfg "github.com/wb-go/wbf/config"
)

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
}

type PostgresConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	Password string        `mapstructure:"password"`
	Database int           `mapstructure:"database"`
	TTL      time.Duration `mapstructure:"ttl"`
}

type ClientConnConfig struct {
	URL     string        `mapstructure:"url"`
	Timeout time.Duration `mapstructure:"timeout"`
}

type RabbitMQConfig struct {
	ClientConn ClientConnConfig `mapstructure:"client_conn"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type GinConfig struct {
	Mode string `mapstructure:"mode"`
}

type TgBotConfig struct {
	Token string `mapstructure:"token"`
}

type EmailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"db"`
	Redis    RedisConfig    `mapstructure:"redis"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Gin      GinConfig      `mapstructure:"gin"`
	TgBot    TgBotConfig    `mapstructure:"tg_bot"`
	Email    EmailConfig    `mapstructure:"email"`
}

func New() (*Config, error) {
	config := wbfcfg.New()

	cfgFile := "./config/config.yaml"
	if err := config.LoadConfigFiles(cfgFile); err != nil {
		return nil, fmt.Errorf("load config files: %w", err)
	}

	config.EnableEnv("")

	var cfg Config
	if err := config.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config unmarshal: %w", err)
	}

	return &cfg, nil
}
