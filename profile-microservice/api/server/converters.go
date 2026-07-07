package server

import (
	profilepb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/profile/langs/go"

	"profile-service/infra/models/domain"
)

func protoToUpsertProfile(req *profilepb.UpsertProfileRequest) *domain.UpsertProfileRequest {
	return &domain.UpsertProfileRequest{
		UserID:    req.GetUserId(),
		FullName:  req.GetFullName(),
		Phone:     req.GetPhone(),
		Apartment: req.GetApartment(),
		HouseID:   req.GetHouseId(),
	}
}

func protoToCreateCompany(req *profilepb.CreateManagementCompanyRequest) *domain.CreateManagementCompanyRequest {
	return &domain.CreateManagementCompanyRequest{Name: req.GetName()}
}

func protoToCreateHouse(req *profilepb.CreateHouseRequest) *domain.CreateHouseRequest {
	return &domain.CreateHouseRequest{
		Name:    req.GetName(),
		Address: req.GetAddress(),
		UKID:    req.GetUkId(),
	}
}

func protoToListHouses(req *profilepb.ListHousesRequest) *domain.ListHousesRequest {
	return &domain.ListHousesRequest{UKID: req.GetUkId()}
}

func domainToProtoProfile(p *domain.Profile) *profilepb.Profile {
	return &profilepb.Profile{
		UserId:    p.UserID,
		FullName:  p.FullName,
		Phone:     p.Phone,
		Apartment: p.Apartment,
		HouseId:   p.HouseID,
		UkId:      p.UKID,
	}
}

func domainToProtoCompany(c *domain.ManagementCompany) *profilepb.ManagementCompany {
	return &profilepb.ManagementCompany{Id: c.ID, Name: c.Name}
}

func domainToProtoHouse(h *domain.House) *profilepb.House {
	return &profilepb.House{
		Id:      h.ID,
		Name:    h.Name,
		Address: h.Address,
		UkId:    h.UKID,
	}
}
