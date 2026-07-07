package api

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"requests-service/resources"
	"testing"
	"time"
)

type brokenListener struct{}

func (brokenListener) Accept() (net.Conn, error) { return nil, errors.New("accept") }
func (brokenListener) Close() error              { return nil }
func (brokenListener) Addr() net.Addr            { return &net.TCPAddr{} }
func TestStart(t *testing.T) {
	l := zerolog.Nop()
	a := NewAPI(&resources.Resources{Env: &resources.Env{Host: "h", Port: 1, RequestTimeout: time.Second}, Logger: &l})
	old := listen
	t.Cleanup(func() { listen = old })
	listen = func(context.Context, string, string) (net.Listener, error) { return nil, errors.New("listen") }
	require.EqualError(t, a.Start(context.Background()), "listen")
	listen = func(context.Context, string, string) (net.Listener, error) { return brokenListener{}, nil }
	assert.EqualError(t, a.Start(context.Background()), "accept")
}
