package api

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"news-service/resources"
)

type brokenListener struct{}

func (brokenListener) Accept() (net.Conn, error) { return nil, errors.New("accept") }
func (brokenListener) Close() error              { return nil }
func (brokenListener) Addr() net.Addr            { return &net.TCPAddr{} }

func TestAPIStart(t *testing.T) {
	logger := zerolog.Nop()
	a := NewAPI(&resources.Resources{Env: &resources.Env{Host: "h", Port: 1, RequestTimeout: time.Second}, Logger: &logger})
	original := listen
	t.Cleanup(func() { listen = original })
	listen = func(context.Context, string, string) (net.Listener, error) { return nil, errors.New("listen") }
	require.EqualError(t, a.Start(context.Background()), "listen")
	listen = func(context.Context, string, string) (net.Listener, error) { return brokenListener{}, nil }
	assert.EqualError(t, a.Start(context.Background()), "accept")
}
