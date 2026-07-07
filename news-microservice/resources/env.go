package resources

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	ServiceName string `envconfig:"APP_NAME" default:"news-service"`

	Host string `envconfig:"APP_HOST" default:"0.0.0.0"`
	Port int    `envconfig:"APP_PORT" default:"50053"`

	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	RequestTimeout time.Duration `envconfig:"REQUEST_TIMEOUT" default:"5s"`

	PostgresHost     string `envconfig:"POSTGRES_HOST"     default:"localhost"`
	PostgresPort     int    `envconfig:"POSTGRES_PORT"     default:"5432"`
	PostgresUser     string `envconfig:"POSTGRES_USER"     default:"postgres"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD" default:"postgres"`
	PostgresDB       string `envconfig:"POSTGRES_DB"       default:"news"`
	PostgresSSLMode  string `envconfig:"POSTGRES_SSLMODE"  default:"disable"`

	RabbitHost     string `envconfig:"RABBIT_HOST"     default:"localhost"`
	RabbitPort     int    `envconfig:"RABBIT_PORT"     default:"5672"`
	RabbitUser     string `envconfig:"RABBIT_USER"     default:"guest"`
	RabbitPassword string `envconfig:"RABBIT_PASSWORD" default:"guest"`
}

func (e *Env) Addr() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

func (e *Env) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		e.PostgresHost, e.PostgresPort,
		e.PostgresUser, e.PostgresPassword,
		e.PostgresDB, e.PostgresSSLMode,
	)
}

func (e *Env) RabbitDSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", e.RabbitUser, e.RabbitPassword, e.RabbitHost, e.RabbitPort)
}

func initEnv() (*Env, error) {
	var e Env
	if err := envconfig.Process("", &e); err != nil {
		return nil, fmt.Errorf("init env: %w", err)
	}

	return &e, nil
}
