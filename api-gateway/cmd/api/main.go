package main

import (
	"context"
	"errors"

	"api-gateway/internal/api"
	"api-gateway/pkg/context_os"
	"api-gateway/resources"
	"github.com/rs/zerolog/log"
)

func main() {
	res, err := resources.InitResources()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init resources")
	}

	ctx := context_os.Context(context.Background(), res.Logger)

	if err := api.New(res).Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		res.Logger.Fatal().Err(err).Msg("server stopped")
	}
}
