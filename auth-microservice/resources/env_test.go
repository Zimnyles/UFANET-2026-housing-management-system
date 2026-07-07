package resources

import (
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jwtinfra "auth-service/infra/jwt"
)

func TestEnv(t *testing.T) {
	e := &Env{Host: "host", Port: 42, PostgresHost: "db", PostgresPort: 5432, PostgresUser: "u", PostgresPassword: "p", PostgresDB: "name", PostgresSSLMode: "disable"}
	assert.Equal(t, "host:42", e.Addr())
	assert.Equal(t, "host=db port=5432 user=u password=p dbname=name sslmode=disable", e.DSN())
	t.Setenv("APP_PORT", "1234")
	got, err := initEnv()
	require.NoError(t, err)
	assert.Equal(t, 1234, got.Port)
	assert.Equal(t, 15*time.Minute, got.JWTAccessTTL)
	t.Setenv("APP_PORT", "bad")
	_, err = initEnv()
	require.Error(t, err)
}

func TestAdaptersAndLogger(t *testing.T) {
	logger := initLogger("svc", "debug")
	assert.Equal(t, zerolog.DebugLevel, logger.GetLevel())
	logger = initLogger("svc", "invalid")
	assert.Equal(t, zerolog.InfoLevel, logger.GetLevel())

	manager := &jwtAdapter{m: jwtinfra.NewManager("secret", time.Minute, time.Hour)}
	access, err := manager.GenerateAccess("u", "admin")
	require.NoError(t, err)
	require.NotEmpty(t, access)
	refresh, err := manager.GenerateRefresh("u", "admin")
	require.NoError(t, err)
	claims, err := manager.ParseRefresh(refresh)
	require.NoError(t, err)
	assert.Equal(t, "u", claims.UserID)
	_, err = manager.ParseRefresh("bad")
	require.Error(t, err)

	h := &hasherAdapter{}
	hash, err := h.Hash("password")
	require.NoError(t, err)
	assert.True(t, h.Check("password", hash))
	assert.False(t, h.Check("bad", hash))
}
