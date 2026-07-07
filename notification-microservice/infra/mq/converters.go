package mq

import "notification-service/infra/models/domain"

func requestEventToNotification(event requestEvent) *domain.Notification {
	return &domain.Notification{UserID: event.UserID, Type: "request_status", Title: "Статус заявки изменён", Body: "Новый статус: " + statusLabel(event.Status)}
}

func newsEventToNotification(event newsEvent) *domain.Notification {
	return &domain.Notification{HouseID: event.HouseID, Type: "news", Title: "Новая новость дома", Body: event.Title}
}

func statusLabel(status string) string {
	switch status {
	case "open":
		return "открыта"
	case "in_progress":
		return "в работе"
	case "done":
		return "выполнена"
	case "cancelled":
		return "отменена"
	default:
		return status
	}
}
