package notifications_client

import (
	"context"

	notificationspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/notifications/langs/go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	ClientConnection          interface{ Close() error }
	NotificationServiceClient interface {
		RegisterDevice(context.Context, *notificationspb.RegisterDeviceRequest, ...grpc.CallOption) (*emptypb.Empty, error)
		UnregisterDevice(context.Context, *notificationspb.UnregisterDeviceRequest, ...grpc.CallOption) (*emptypb.Empty, error)
		Subscribe(context.Context, *notificationspb.SubscribeRequest, ...grpc.CallOption) (notificationspb.NotificationService_SubscribeClient, error)
	}
)
