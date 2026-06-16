package server

import (
	"profile-service/infra/models/dto"

	profilepb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/profile/langs/go"
)

// ─── inbound: proto → dto ─────────────────────────────────────────────────────

func protoToUpsertProfileDTO(req *profilepb.UpsertProfileRequest) *dto.UpsertProfileRequest {
	return &dto.UpsertProfileRequest{
		UserID:    req.GetUserId(),
		FullName:  req.GetFullName(),
		Phone:     req.GetPhone(),
		Apartment: req.GetApartment(),
		HouseID:   req.GetHouseId(),
	}
}

func protoToCreateCompanyDTO(req *profilepb.CreateManagementCompanyRequest) *dto.CreateManagementCompanyRequest {
	return &dto.CreateManagementCompanyRequest{Name: req.GetName()}
}

func protoToCreateHouseDTO(req *profilepb.CreateHouseRequest) *dto.CreateHouseRequest {
	return &dto.CreateHouseRequest{
		Name:    req.GetName(),
		Address: req.GetAddress(),
		UKID:    req.GetUkId(),
	}
}

func protoToListHousesDTO(req *profilepb.ListHousesRequest) *dto.ListHousesRequest {
	return &dto.ListHousesRequest{UKID: req.GetUkId()}
}

// ─── outbound: dto → proto ────────────────────────────────────────────────────

func dtoToProtoProfile(p *dto.Profile) *profilepb.Profile {
	return &profilepb.Profile{
		UserId:    p.UserID,
		FullName:  p.FullName,
		Phone:     p.Phone,
		Apartment: p.Apartment,
		HouseId:   p.HouseID,
		UkId:      p.UKID,
	}
}

func dtoToProtoCompany(c *dto.ManagementCompany) *profilepb.ManagementCompany {
	return &profilepb.ManagementCompany{Id: c.ID, Name: c.Name}
}

func dtoToProtoHouse(h *dto.House) *profilepb.House {
	return &profilepb.House{
		Id:      h.ID,
		Name:    h.Name,
		Address: h.Address,
		UkId:    h.UKID,
	}
}
