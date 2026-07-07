package jwt

import (
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {
	m := NewManager("secret", time.Minute, time.Hour)
	access, err := m.GenerateAccess("u1", "admin")
	require.NoError(t, err)
	require.NotEmpty(t, access)
	refresh, err := m.GenerateRefresh("u1", "admin")
	require.NoError(t, err)
	claims, err := m.ParseRefresh(refresh)
	require.NoError(t, err)
	assert.Equal(t, "u1", claims.UserID)
	assert.Equal(t, "admin", claims.Role)

	_, err = m.ParseRefresh("invalid")
	require.Error(t, err)
	other := NewManager("other", time.Minute, time.Hour)
	_, err = other.ParseRefresh(refresh)
	require.Error(t, err)

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodNone, &Claims{})
	none, err := token.SignedString(jwtlib.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	_, err = m.ParseRefresh(none)
	require.Error(t, err)

	expired := NewManager("secret", -time.Hour, -time.Hour)
	expiredToken, err := expired.GenerateRefresh("u", "user")
	require.NoError(t, err)
	_, err = m.ParseRefresh(expiredToken)
	require.Error(t, err)
}
