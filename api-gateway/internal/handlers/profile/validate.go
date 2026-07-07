package profile_handler

import (
	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/dto"
)

func validateUpsertProfile(req *dto.UpsertProfileRequest) error {
	if req.FullName == "" {
		return app_errors.ErrProfileFullNameRequired
	}

	if len(req.FullName) < 2 {
		return app_errors.ErrProfileFullNameTooShort
	}

	if len(req.FullName) > 100 {
		return app_errors.ErrProfileFullNameTooLong
	}

	if req.Phone == "" {
		return app_errors.ErrProfilePhoneRequired
	}

	if len(req.Phone) > 20 {
		return app_errors.ErrProfilePhoneTooLong
	}

	if req.Apartment == "" {
		return app_errors.ErrProfileApartmentRequired
	}

	if len(req.Apartment) > 20 {
		return app_errors.ErrProfileApartmentTooLong
	}

	if req.HouseID == "" {
		return app_errors.ErrProfileHouseIDRequired
	}

	return nil
}

func validateCreateCompany(req *dto.CreateManagementCompanyRequest) error {
	if req.Name == "" {
		return app_errors.ErrCompanyNameRequired
	}

	if len(req.Name) > 200 {
		return app_errors.ErrCompanyNameTooLong
	}

	return nil
}

func validateCreateHouse(req *dto.CreateHouseRequest) error {
	if req.Name == "" {
		return app_errors.ErrHouseNameRequired
	}

	if req.Address == "" {
		return app_errors.ErrHouseAddressRequired
	}

	if req.UKID == "" {
		return app_errors.ErrHouseUKIDRequired
	}

	return nil
}
