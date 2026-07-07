package notifications_handler

import (
	"api-gateway/internal/models/domain"
	"api-gateway/internal/models/dto"
)

func toDomainRegister(req dto.RegisterDeviceRequest) domain.RegisterDevice {
	return domain.RegisterDevice{DeviceToken: req.DeviceToken, Platform: req.Platform}
}
