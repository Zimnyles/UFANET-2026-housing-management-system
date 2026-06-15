package dto

import "api-gateway/internal/models/domain"

type RegisterDeviceRequest struct {
	DeviceToken string          `json:"device_token"`
	Platform    domain.Platform `json:"platform"`
}
