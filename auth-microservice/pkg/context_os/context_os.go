package context_os

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
)

var defaultSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

func Context(ctx context.Context, logger *zerolog.Logger) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, defaultSignals...)

	go func() {
		defer signal.Stop(sigs)

		select {
		case <-ctx.Done():
		case sig := <-sigs:
			logger.Info().Msgf("got signal %s", sig)
			cancel()
		}
	}()

	return ctx
}
