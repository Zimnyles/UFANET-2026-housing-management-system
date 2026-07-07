package main

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"

	"profile-service/api"
	"profile-service/pkg/context_os"
	"profile-service/resources"
)

func main() {
	ctx := context.Background()

	res, err := resources.InitResources(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init resources")
	}

	ctx = context_os.Context(ctx, res.Logger)

	if err := api.NewAPI(res).Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		res.Logger.Fatal().Err(err).Msg("server stopped")
	}
}
