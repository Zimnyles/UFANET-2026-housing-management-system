package resources

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	ServiceName string `envconfig:"APP_NAME" default:"auth-service"`

	Host string `envconfig:"APP_HOST" default:"0.0.0.0"`
	Port int    `envconfig:"APP_PORT" default:"50051"`

	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	JWTSecret     string        `envconfig:"JWT_SECRET"      default:"dev-secret-change-in-production"`
	JWTAccessTTL  time.Duration `envconfig:"JWT_ACCESS_TTL"  default:"15m"`
	JWTRefreshTTL time.Duration `envconfig:"JWT_REFRESH_TTL" default:"168h"`

	RequestTimeout time.Duration `envconfig:"REQUEST_TIMEOUT" default:"5s"`

	PostgresHost     string `envconfig:"POSTGRES_HOST"     default:"localhost"`
	PostgresPort     int    `envconfig:"POSTGRES_PORT"     default:"5432"`
	PostgresUser     string `envconfig:"POSTGRES_USER"     default:"postgres"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD" default:"postgres"`
	PostgresDB       string `envconfig:"POSTGRES_DB"       default:"auth"`
	PostgresSSLMode  string `envconfig:"POSTGRES_SSLMODE"  default:"disable"`
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

func initEnv() (*Env, error) {
	var e Env
	if err := envconfig.Process("", &e); err != nil {
		return nil, fmt.Errorf("init env: %w", err)
	}

	return &e, nil
}
