package config

import (
	"time"

	"github.com/wb-go/wbf/config"
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

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type WebConfig struct {
	Path string `mapstructure:"path"`
}

type RedisConfig struct {
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	Password string        `mapstructure:"password"`
	Database int           `mapstructure:"database"`
	TTL      time.Duration `mapstructure:"ttl"`
}

type RabbitMQConnConfig struct {
	URL     string        `mapstructure:"url"`
	Retries int           `mapstructure:"retries"`
	Pause   time.Duration `mapstructure:"pause"`
}

type RabbitMQExchangeConfig struct {
	Name       string `mapstructure:"name"`
	Type       string `mapstructure:"type"`
	Durable    bool   `mapstructure:"durable"`
	AutoDelete bool   `mapstructure:"auto_delete"`
	Internal   bool   `mapstructure:"internal"`
	NoWait     bool   `mapstructure:"no_wait"`
}

type RabbitMQQueueConfig struct {
	Name       string `mapstructure:"name"`
	RoutingKey string `mapstructure:"routing_key"`
	Durable    bool   `mapstructure:"durable"`
	AutoDelete bool   `mapstructure:"auto_delete"`
	Exclusive  bool   `mapstructure:"exclusive"`
	NoWait     bool   `mapstructure:"no_wait"`
}

type RabbitMQConfig struct {
	Conn            RabbitMQConnConfig     `mapstructure:"conn"`
	Exchange        RabbitMQExchangeConfig `mapstructure:"exchange"`
	RetryExchange   RabbitMQExchangeConfig `mapstructure:"retry_exchange"`
	QueueEmail      RabbitMQQueueConfig    `mapstructure:"queue_email"`
	QueueTg         RabbitMQQueueConfig    `mapstructure:"queue_tg"`
	RetryQueueEmail RabbitMQQueueConfig    `mapstructure:"retry_queue_email"`
	RetryQueueTg    RabbitMQQueueConfig    `mapstructure:"retry_queue_tg"`
	Attempts        int                    `mapstructure:"attempts"`
	Delay           time.Duration          `mapstructure:"delay"`
	Backoff         float64                `mapstructure:"backoff"`
}

type EmailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
	SSL      bool   `mapstructure:"ssl"`
}

type TelegramConfig struct {
	BotToken string        `mapstructure:"bot_token"`
	Timeout  time.Duration `mapstructure:"timeout"`
	Debug    bool          `mapstructure:"debug"`
}

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	DB       PostgresConfig `mapstructure:"postgres"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Web      WebConfig      `mapstructure:"web"`
	Redis    RedisConfig    `mapstructure:"redis"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	Email    EmailConfig    `mapstructure:"email"`
	Telegram TelegramConfig `mapstructure:"telegram"`
}

func New() (*Config, error) {
	cfg := config.New()
	cfg.LoadConfigFiles("./config/config.yaml")

	// Включить env переменные с приставкой
	cfg.EnableEnv("")

	var config Config
	err := cfg.Unmarshal(&config)
	return &config, err
}
