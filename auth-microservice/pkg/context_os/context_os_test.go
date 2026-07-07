package context_os

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestContextFollowsParentCancellation(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	logger := zerolog.Nop()
	ctx := Context(parent, &logger)
	cancel()
	<-ctx.Done()
	require.ErrorIs(t, ctx.Err(), context.Canceled)
}
