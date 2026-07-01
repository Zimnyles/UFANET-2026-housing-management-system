package api

import (
	"context"
	"net"

	"news-service/api/interceptors"
	"news-service/api/server"
	"news-service/resources"
	news_service "news-service/services/news"

	newspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/news/langs/go"
	"google.golang.org/grpc"
)

type API struct {
	res *resources.Resources
}

func NewAPI(res *resources.Resources) *API {
	return &API{res: res}
}

func (a *API) Start(ctx context.Context) error {
	svc := news_service.New(a.res.Repo, a.res.Publisher, a.res.Logger)
	newsServer := server.NewNewsServer(svc, a.res.Logger)

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.Timeout(a.res.Env.RequestTimeout),
			interceptors.ErrorMapping(),
			interceptors.Validation(),
		),
	)

	newspb.RegisterNewsServiceServer(srv, newsServer)

	lis, err := net.Listen("tcp", a.res.Env.Addr())
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
