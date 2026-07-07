package api

import (
	"context"
	"net"

	profilepb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/profile/langs/go"
	"google.golang.org/grpc"

	"profile-service/api/interceptors"
	"profile-service/api/server"
	"profile-service/resources"
	profile_service "profile-service/services/profile"
)

var listen = func(ctx context.Context, network, address string) (net.Listener, error) {
	var c net.ListenConfig

	return c.Listen(ctx, network, address)
}

type API struct {
	res *resources.Resources
}

func NewAPI(res *resources.Resources) *API {
	return &API{res: res}
}

func (a *API) Start(ctx context.Context) error {
	svc := profile_service.New(a.res.Repo, a.res.Logger)
	profileServer := server.NewProfileServer(svc, a.res.Logger)

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.Timeout(a.res.Env.RequestTimeout),
			interceptors.ErrorMapping(),
			interceptors.Validation(),
		),
	)

	profilepb.RegisterProfileServiceServer(srv, profileServer)

	lis, err := listen(ctx, "tcp", a.res.Env.Addr())
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
