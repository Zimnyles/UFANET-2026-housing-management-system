package resources

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"requests-service/infra/mq"
)

func TestEnv(t *testing.T) {
	e := &Env{Host: "h", Port: 1, PostgresHost: "d", PostgresPort: 2, PostgresUser: "u", PostgresPassword: "p", PostgresDB: "n", PostgresSSLMode: "disable", RabbitHost: "r", RabbitPort: 3, RabbitUser: "x", RabbitPassword: "y"}
	assert.Equal(t, "h:1", e.Addr())
	assert.Equal(t, "host=d port=2 user=u password=p dbname=n sslmode=disable", e.DSN())
	assert.Equal(t, "amqp://x:y@r:3/", e.RabbitDSN())
	t.Setenv("APP_PORT", "42")
	v, err := initEnv()
	require.NoError(t, err)
	assert.Equal(t, 42, v.Port)
	t.Setenv("APP_PORT", "bad")
	_, err = initEnv()
	require.Error(t, err)
}

func TestResources(t *testing.T) {
	assert.Equal(t, zerolog.DebugLevel, initLogger("s", "debug").GetLevel())
	assert.Equal(t, zerolog.InfoLevel, initLogger("s", "bad").GetLevel())
	l := zerolog.Nop()
	(&Resources{Logger: &l}).Close()
	(&Resources{Logger: &l, Publisher: &mq.Publisher{}}).Close()
	t.Setenv("APP_PORT", "bad")
	_, e := InitResources(context.Background())
	require.Error(t, e)
	t.Setenv("APP_PORT", "1")
	t.Setenv("POSTGRES_HOST", "127.0.0.1")
	t.Setenv("POSTGRES_PORT", "1")
	ctx, c := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer c()
	_, e = InitResources(ctx)
	require.Error(t, e)
}
