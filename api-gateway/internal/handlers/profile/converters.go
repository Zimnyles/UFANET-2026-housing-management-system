package profile_handler

import (
	"api-gateway/internal/models/domain"
	"api-gateway/internal/models/dto"
)

func toDomainUpsertRequest(userID string, req dto.UpsertProfileRequest) *domain.UpsertProfileRequest {
	return &domain.UpsertProfileRequest{
		UserID:    userID,
		FullName:  req.FullName,
		Phone:     req.Phone,
		Apartment: req.Apartment,
		HouseID:   req.HouseID,
	}
}

func toDTOProfileResponse(p *domain.Profile) dto.ProfileResponse {
	return dto.ProfileResponse{
		UserID:    p.UserID,
		FullName:  p.FullName,
		Phone:     p.Phone,
		Apartment: p.Apartment,
		HouseID:   p.HouseID,
		UKID:      p.UKID,
	}
}

func toDomainCreateCompanyRequest(req dto.CreateManagementCompanyRequest) *domain.CreateManagementCompanyRequest {
	return &domain.CreateManagementCompanyRequest{Name: req.Name}
}

func toDTOCompanyResponse(c *domain.ManagementCompany) dto.ManagementCompanyResponse {
	return dto.ManagementCompanyResponse{ID: c.ID, Name: c.Name}
}

func toDomainCreateHouseRequest(req dto.CreateHouseRequest) *domain.CreateHouseRequest {
	return &domain.CreateHouseRequest{Name: req.Name, Address: req.Address, UKID: req.UKID}
}

func toDTOHouseResponse(h *domain.House) dto.HouseResponse {
	return dto.HouseResponse{ID: h.ID, Name: h.Name, Address: h.Address, UKID: h.UKID}
}
