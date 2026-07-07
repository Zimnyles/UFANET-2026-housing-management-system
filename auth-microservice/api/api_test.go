package api

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"auth-service/resources"
)

func TestAPIStart(t *testing.T) {
	logger := zerolog.Nop()
	t.Run("listen error", func(t *testing.T) {
		a := NewAPI(&resources.Resources{Env: &resources.Env{Host: "invalid host", Port: -1, RequestTimeout: time.Second}, Logger: &logger})
		err := a.Start(context.Background())
		require.Error(t, err)
	})

	t.Run("graceful shutdown", func(t *testing.T) {
		a := NewAPI(&resources.Resources{Env: &resources.Env{Host: "127.0.0.1", Port: 0, RequestTimeout: time.Second}, Logger: &logger})
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if err := a.Start(ctx); err != nil {
			require.True(t, strings.Contains(err.Error(), "listen") || strings.Contains(err.Error(), "bind"))
		}
	})
}
