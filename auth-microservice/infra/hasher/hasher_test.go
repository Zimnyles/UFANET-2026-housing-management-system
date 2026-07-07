package hasher

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashAndCheck(t *testing.T) {
	hash, err := Hash("password")
	require.NoError(t, err)
	assert.True(t, Check("password", hash))
	assert.False(t, Check("wrong", hash))
	_, err = Hash(strings.Repeat("x", 73))
	require.Error(t, err)
}
