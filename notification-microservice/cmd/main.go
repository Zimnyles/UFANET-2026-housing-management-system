package main

import (
	"context"
	"errors"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"notification-service/api"
	"notification-service/resources"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	res, err := resources.InitResources(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("init resources")
	}

	defer cancel()
	defer res.Close()

	if err := api.New(res).Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		res.Logger.Fatal().Err(err).Msg("server stopped")
	}
}
