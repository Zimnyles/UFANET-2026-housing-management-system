package context_os

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

var defaultSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

func Context(ctx context.Context, logger *zerolog.Logger) context.Context {
	return getContext(ctx, 0, logger)
}

func getContext(ctx context.Context, timeout time.Duration, logger *zerolog.Logger) context.Context {
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

			if timeout == 0 {
				return
			}

			logger.Info().Msgf("wait for timeout %s", timeout)
			time.Sleep(timeout)
			logger.Info().Msgf("timeout %s reached, exiting...", timeout)

			os.Exit(1)
		}
	}()

	return ctx
}
