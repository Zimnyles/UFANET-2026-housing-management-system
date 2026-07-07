package server

import (
	"time"

	notificationspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/notifications/langs/go"

	"notification-service/infra/models/domain"
)

func protoToDevice(req *notificationspb.RegisterDeviceRequest) *domain.Device {
	return &domain.Device{UserID: req.GetUserId(), Token: req.GetDeviceToken(), Platform: req.GetPlatform()}
}

func protoToList(req *notificationspb.GetUserNotificationsRequest) *domain.ListRequest {
	return &domain.ListRequest{UserID: req.GetUserId(), HouseID: req.GetHouseId(), Limit: int(req.GetLimit()), Offset: int(req.GetOffset())}
}

func domainToProto(n *domain.Notification) *notificationspb.Notification {
	return &notificationspb.Notification{Id: n.ID, UserId: n.UserID, HouseId: n.HouseID, Type: n.Type, Title: n.Title, Body: n.Body, CreatedAt: n.CreatedAt.UTC().Format(time.RFC3339), Read: n.Read}
}

func domainListToProto(items []*domain.Notification) []*notificationspb.Notification {
	result := make([]*notificationspb.Notification, 0, len(items))
	for _, item := range items {
		result = append(result, domainToProto(item))
	}

	return result
}
