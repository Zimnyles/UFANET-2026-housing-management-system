package repository

import "notification-service/infra/models/domain"

func deviceToRow(device *domain.Device) *dbDevice {
	return &dbDevice{UserID: device.UserID, DeviceToken: device.Token, Platform: device.Platform}
}

func notificationToRow(n *domain.Notification) *dbNotification {
	return &dbNotification{UserID: optionalString(n.UserID), HouseID: optionalString(n.HouseID), Type: n.Type, Title: n.Title, Body: n.Body, CreatedAt: n.CreatedAt, Read: n.Read}
}

func notificationToDomain(row *dbNotification) *domain.Notification {
	return &domain.Notification{ID: row.ID, UserID: stringValue(row.UserID), HouseID: stringValue(row.HouseID), Type: row.Type, Title: row.Title, Body: row.Body, CreatedAt: row.CreatedAt, Read: row.Read}
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
