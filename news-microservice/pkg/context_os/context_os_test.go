package context_os

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := zerolog.Nop()
	ctx := Context(context.Background(), &logger)
	assert.Same(t, &logger, Logger(ctx))
	assert.Nil(t, Logger(context.Background()))
}
