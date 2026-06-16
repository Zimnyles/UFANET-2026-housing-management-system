package server

import (
	infra_errors "profile-service/infra/errors"

	profilepb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/profile/langs/go"
)

func ValidateGetProfileRequest(req *profilepb.GetProfileRequest) error {
	if req.GetUserId() == "" {
		return infra_errors.ErrUserIDRequired
	}
	return nil
}

func ValidateUpsertProfileRequest(req *profilepb.UpsertProfileRequest) error {
	if req.GetUserId() == "" {
		return infra_errors.ErrUserIDRequired
	}
	if req.GetFullName() == "" {
		return infra_errors.ErrFullNameRequired
	}
	if len(req.GetFullName()) < 2 {
		return infra_errors.ErrFullNameTooShort
	}
	if len(req.GetFullName()) > 100 {
		return infra_errors.ErrFullNameTooLong
	}
	if req.GetPhone() == "" {
		return infra_errors.ErrPhoneRequired
	}
	if len(req.GetPhone()) > 20 {
		return infra_errors.ErrPhoneTooLong
	}
	if req.GetApartment() == "" {
		return infra_errors.ErrApartmentRequired
	}
	if len(req.GetApartment()) > 20 {
		return infra_errors.ErrApartmentTooLong
	}
	if req.GetHouseId() == "" {
		return infra_errors.ErrHouseIDRequired
	}
	return nil
}

func ValidateCreateManagementCompanyRequest(req *profilepb.CreateManagementCompanyRequest) error {
	if req.GetName() == "" {
		return infra_errors.ErrCompanyNameRequired
	}
	if len(req.GetName()) > 200 {
		return infra_errors.ErrCompanyNameTooLong
	}
	return nil
}

func ValidateCreateHouseRequest(req *profilepb.CreateHouseRequest) error {
	if req.GetName() == "" {
		return infra_errors.ErrHouseNameRequired
	}
	if req.GetAddress() == "" {
		return infra_errors.ErrHouseAddressRequired
	}
	if req.GetUkId() == "" {
		return infra_errors.ErrUKIDRequired
	}
	return nil
}
