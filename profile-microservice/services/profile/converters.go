package profile_service

import "profile-service/infra/models/domain"

func upsertRequestToProfile(req *domain.UpsertProfileRequest) *domain.Profile {
	return &domain.Profile{
		UserID: req.UserID, FullName: req.FullName, Phone: req.Phone,
		Apartment: req.Apartment, HouseID: req.HouseID,
	}
}

func createCompanyRequestToModel(req *domain.CreateManagementCompanyRequest) *domain.ManagementCompany {
	return &domain.ManagementCompany{Name: req.Name}
}

func createHouseRequestToModel(req *domain.CreateHouseRequest) *domain.House {
	return &domain.House{Name: req.Name, Address: req.Address, UKID: req.UKID}
}
