package server

import (
	"context"
	"io"
	"math"

	notificationspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/notifications/langs/go"
	"google.golang.org/protobuf/types/known/emptypb"

	infra_errors "notification-service/infra/errors"
	"notification-service/infra/grpcerrors"
)

type Server struct {
	notificationspb.UnimplementedNotificationServiceServer
	service NotificationsService
}

func New(service NotificationsService) *Server { return &Server{service: service} }

func (s *Server) RegisterDevice(ctx context.Context, req *notificationspb.RegisterDeviceRequest) (*emptypb.Empty, error) {
	if req.GetUserId() == "" || req.GetDeviceToken() == "" || req.GetPlatform() == "" {
		return nil, grpcerrors.ToGRPCError(validateRegisterRequest(req))
	}

	if err := s.service.Register(ctx, protoToDevice(req)); err != nil {
		return nil, grpcerrors.ToGRPCError(infra_errors.ErrRegisterDevice)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) UnregisterDevice(ctx context.Context, req *notificationspb.UnregisterDeviceRequest) (*emptypb.Empty, error) {
	if err := s.service.Unregister(ctx, req.GetUserId(), req.GetDeviceToken()); err != nil {
		return nil, grpcerrors.ToGRPCError(infra_errors.ErrUnregisterDevice)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) GetUserNotifications(ctx context.Context, req *notificationspb.GetUserNotificationsRequest) (*notificationspb.GetUserNotificationsResponse, error) {
	items, total, err := s.service.List(ctx, protoToList(req))
	if err != nil {
		return nil, grpcerrors.ToGRPCError(infra_errors.ErrListNotifications)
	}

	return &notificationspb.GetUserNotificationsResponse{Notifications: domainListToProto(items), Total: totalToInt32(total)}, nil
}

func totalToInt32(total int64) int32 {
	if total <= 0 {
		return 0
	}

	if total > math.MaxInt32 {
		return math.MaxInt32
	}

	return int32(total)
}

func (s *Server) Subscribe(req *notificationspb.SubscribeRequest, stream notificationspb.NotificationService_SubscribeServer) error {
	if req.GetUserId() == "" {
		return grpcerrors.ToGRPCError(infra_errors.ErrUserIDRequired)
	}

	events, unsubscribe := s.service.Subscribe(req.GetUserId(), req.GetHouseId())
	defer unsubscribe()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case event, ok := <-events:
			if !ok {
				return io.EOF
			}

			if err := stream.Send(domainToProto(event)); err != nil {
				return err
			}
		}
	}
}

func validateRegisterRequest(req *notificationspb.RegisterDeviceRequest) error {
	if req.GetUserId() == "" {
		return infra_errors.ErrUserIDRequired
	}

	if req.GetDeviceToken() == "" {
		return infra_errors.ErrDeviceTokenRequired
	}

	if req.GetPlatform() == "" {
		return infra_errors.ErrPlatformRequired
	}

	return nil
}
