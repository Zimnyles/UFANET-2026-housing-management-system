package api

import (
	"context"
	"net"

	requestspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/requests/langs/go"
	"google.golang.org/grpc"

	"requests-service/api/interceptors"
	"requests-service/api/server"
	"requests-service/resources"
	requests_service "requests-service/services/requests"
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
	svc := requests_service.New(a.res.Repo, a.res.Publisher, a.res.Logger)
	requestsServer := server.NewRequestsServer(svc, a.res.Logger)

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.Timeout(a.res.Env.RequestTimeout),
			interceptors.ErrorMapping(),
			interceptors.Validation(),
		),
	)

	requestspb.RegisterRequestsServiceServer(srv, requestsServer)

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
