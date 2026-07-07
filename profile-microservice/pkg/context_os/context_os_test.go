package context_os

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogger(t *testing.T) {
	l := zerolog.Nop()
	assert.Same(t, &l, Logger(Context(context.Background(), &l)))
	assert.Nil(t, Logger(context.Background()))
}
