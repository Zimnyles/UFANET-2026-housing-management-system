package resources

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestEnv(t *testing.T) {
	e := &Env{Host: "h", Port: 1, PostgresHost: "db", PostgresPort: 2, PostgresUser: "u", PostgresPassword: "p", PostgresDB: "n", PostgresSSLMode: "disable"}
	assert.Equal(t, "h:1", e.Addr())
	assert.Equal(t, "host=db port=2 user=u password=p dbname=n sslmode=disable", e.DSN())
	t.Setenv("APP_PORT", "42")
	v, err := initEnv()
	require.NoError(t, err)
	assert.Equal(t, 42, v.Port)
	t.Setenv("APP_PORT", "bad")
	_, err = initEnv()
	require.Error(t, err)
}
func TestResourcesErrorsAndLogger(t *testing.T) {
	assert.Equal(t, zerolog.DebugLevel, initLogger("s", "debug").GetLevel())
	assert.Equal(t, zerolog.InfoLevel, initLogger("s", "bad").GetLevel())
	t.Setenv("APP_PORT", "bad")
	_, e := InitResources(context.Background())
	require.Error(t, e)
	t.Setenv("APP_PORT", "1")
	t.Setenv("POSTGRES_HOST", "127.0.0.1")
	t.Setenv("POSTGRES_PORT", "1")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_, e = InitResources(ctx)
	require.Error(t, e)
}
