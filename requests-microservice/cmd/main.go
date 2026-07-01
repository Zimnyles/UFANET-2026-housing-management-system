package main

import (
	"context"
	"errors"
	"requests-service/api"
	"requests-service/pkg/context_os"
	"requests-service/resources"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	res, err := resources.InitResources(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init resources")
	}
	defer res.Close()

	ctx = context_os.Context(ctx, res.Logger)

	if err := api.NewAPI(res).Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		res.Logger.Fatal().Err(err).Msg("server stopped")
	}
}
