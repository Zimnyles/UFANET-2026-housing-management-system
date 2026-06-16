package resources

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	ServiceName string `envconfig:"APP_NAME"    default:"api-gateway"`

	Host string `envconfig:"APP_HOST"    default:"0.0.0.0"`
	Port int    `envconfig:"APP_PORT"    default:"8080"`

	LogLevel string `envconfig:"LOG_LEVEL"   default:"info"`

	JWTSecret string `envconfig:"JWT_SECRET"  default:"dev-secret-change-in-production"`

	RateLimitMax        int           `envconfig:"RATE_LIMIT_MAX"        default:"10"`
	RateLimitExpiration time.Duration `envconfig:"RATE_LIMIT_EXPIRATION" default:"1m"`

	Prefork bool `envconfig:"APP_PREFORK" default:"false"`

	RequestTimeout time.Duration `envconfig:"REQUEST_TIMEOUT" default:"200ms"`

	AuthAddr string `envconfig:"AUTH_ADDR" default:"localhost:50051"`

	RedisHost     string `envconfig:"REDIS_HOST"     default:"localhost"`
	RedisPort     int    `envconfig:"REDIS_PORT"     default:"6379"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`
	RedisDB       int    `envconfig:"REDIS_DB"       default:"0"`
}

func (e *Env) Addr() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

func initEnv() (*Env, error) {
	var e Env
	if err := envconfig.Process("", &e); err != nil {
		return nil, fmt.Errorf("failed to init env: %w", err)
	}

	return &e, nil
}
