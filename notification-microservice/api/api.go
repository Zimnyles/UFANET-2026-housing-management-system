package api

import (
	"context"
	"net"

	notificationspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/notifications/langs/go"
	"google.golang.org/grpc"

	"notification-service/api/server"
	"notification-service/infra/mq"
	"notification-service/resources"
	notifications_service "notification-service/services/notifications"
)

type API struct{ res *resources.Resources }

func New(res *resources.Resources) *API { return &API{res: res} }
func (a *API) Start(ctx context.Context) error {
	service := notifications_service.New(a.res.Repo, a.res.Logger)

	consumer := mq.New(a.res.Env.RabbitDSN(), service, a.res.Logger)
	go consumer.Run(ctx)

	grpcServer := grpc.NewServer()
	notificationspb.RegisterNotificationServiceServer(grpcServer, server.New(service))

	var listenConfig net.ListenConfig

	lis, err := listenConfig.Listen(ctx, "tcp", a.res.Env.Addr())
	if err != nil {
		return err
	}

	errors := make(chan error, 1)

	go func() { errors <- grpcServer.Serve(lis) }()

	select {
	case <-ctx.Done():
		grpcServer.GracefulStop()

		return nil
	case err := <-errors:
		return err
	}
}
