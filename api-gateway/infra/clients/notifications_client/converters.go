package notifications_client

import (
	notificationspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/notifications/langs/go"

	"api-gateway/internal/models/domain"
)

func toProtoRegister(userID string, req domain.RegisterDevice) *notificationspb.RegisterDeviceRequest {
	return &notificationspb.RegisterDeviceRequest{UserId: userID, DeviceToken: req.DeviceToken, Platform: string(req.Platform)}
}

func toProtoUnregister(userID, token string) *notificationspb.UnregisterDeviceRequest {
	return &notificationspb.UnregisterDeviceRequest{UserId: userID, DeviceToken: token}
}

func toProtoSubscribe(userID, houseID string) *notificationspb.SubscribeRequest {
	return &notificationspb.SubscribeRequest{UserId: userID, HouseId: houseID}
}

func toDomainNotification(pb *notificationspb.Notification) domain.BrowserNotification {
	url := "/notifications"
	if pb.GetType() == "news" {
		url = "/news"
	}

	if pb.GetType() == "request_status" {
		url = "/requests"
	}

	return domain.BrowserNotification{Type: pb.GetType(), Title: pb.GetTitle(), Body: pb.GetBody(), URL: url, CreatedAt: parseTime(pb.GetCreatedAt())}
}
