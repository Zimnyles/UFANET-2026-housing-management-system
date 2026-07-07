package profile_client

import (
	profilepb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/profile/langs/go"

	"api-gateway/internal/models/domain"
)

func toDomainProfile(pb *profilepb.Profile) *domain.Profile {
	if pb == nil {
		return nil
	}

	return &domain.Profile{
		UserID:    pb.GetUserId(),
		FullName:  pb.GetFullName(),
		Phone:     pb.GetPhone(),
		Apartment: pb.GetApartment(),
		HouseID:   pb.GetHouseId(),
		UKID:      pb.GetUkId(),
	}
}

func toDomainManagementCompany(pb *profilepb.ManagementCompany) *domain.ManagementCompany {
	if pb == nil {
		return nil
	}

	return &domain.ManagementCompany{ID: pb.GetId(), Name: pb.GetName()}
}

func toDomainHouse(pb *profilepb.House) *domain.House {
	if pb == nil {
		return nil
	}

	return &domain.House{
		ID:      pb.GetId(),
		Name:    pb.GetName(),
		Address: pb.GetAddress(),
		UKID:    pb.GetUkId(),
	}
}

func toProtoGetProfileRequest(userID string) *profilepb.GetProfileRequest {
	return &profilepb.GetProfileRequest{UserId: userID}
}

func toProtoUpsertProfileRequest(req *domain.UpsertProfileRequest) *profilepb.UpsertProfileRequest {
	return &profilepb.UpsertProfileRequest{
		UserId:    req.UserID,
		FullName:  req.FullName,
		Phone:     req.Phone,
		Apartment: req.Apartment,
		HouseId:   req.HouseID,
	}
}

func toProtoIsProfileCompleteRequest(userID string) *profilepb.IsProfileCompleteRequest {
	return &profilepb.IsProfileCompleteRequest{UserId: userID}
}

func toProtoCreateManagementCompanyRequest(req *domain.CreateManagementCompanyRequest) *profilepb.CreateManagementCompanyRequest {
	return &profilepb.CreateManagementCompanyRequest{Name: req.Name}
}

func toProtoListManagementCompaniesRequest() *profilepb.ListManagementCompaniesRequest {
	return &profilepb.ListManagementCompaniesRequest{}
}

func toProtoCreateHouseRequest(req *domain.CreateHouseRequest) *profilepb.CreateHouseRequest {
	return &profilepb.CreateHouseRequest{
		Name:    req.Name,
		Address: req.Address,
		UkId:    req.UKID,
	}
}

func toProtoListHousesRequest(ukID string) *profilepb.ListHousesRequest {
	return &profilepb.ListHousesRequest{UkId: ukID}
}
