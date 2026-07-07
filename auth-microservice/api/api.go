package api

import (
	"context"
	"net"

	authpb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/auth/langs/go"
	"google.golang.org/grpc"

	"auth-service/api/interceptors"
	"auth-service/api/server"
	"auth-service/resources"
	authservice "auth-service/services/auth"
)

type API struct {
	res *resources.Resources
}

func NewAPI(res *resources.Resources) *API {
	return &API{res: res}
}

func (a *API) Start(ctx context.Context) error {
	svc := authservice.New(
		a.res.Repo,
		a.res.JWT,
		a.res.Hasher,
		a.res.Logger,
	)

	authServer := server.NewAuthServer(svc, a.res.Logger)

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.Timeout(a.res.Env.RequestTimeout),
			interceptors.ErrorMapping(),
			interceptors.Validation(),
		),
	)

	authpb.RegisterAuthServiceServer(srv, authServer)

	var listenConfig net.ListenConfig

	lis, err := listenConfig.Listen(ctx, "tcp", a.res.Env.Addr())
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)

	go func() {
		a.res.Logger.Info().Str("addr", a.res.Env.Addr()).Msg("starting gRPC server")

		errCh <- srv.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		a.res.Logger.Info().Msg("shutting down gRPC server")
		srv.GracefulStop()

		return nil
	case err := <-errCh:
		return err
	}
}
