package context_os

import (
	"context"

	"github.com/rs/zerolog"
)

type loggerKey struct{}

func Context(ctx context.Context, logger *zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func Logger(ctx context.Context) *zerolog.Logger {
	if l, ok := ctx.Value(loggerKey{}).(*zerolog.Logger); ok {
		return l
	}

	return nil
}
