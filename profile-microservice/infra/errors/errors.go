package errors

import "errors"

var (
	ErrProfileNotFound = errors.New("profile not found")
	ErrInternal        = errors.New("internal server error")
)

var (
	ErrUserIDRequired      = errors.New("user_id is required")
	ErrFullNameRequired    = errors.New("full_name is required")
	ErrFullNameTooShort    = errors.New("full_name min length 2")
	ErrFullNameTooLong     = errors.New("full_name max length 100")
	ErrPhoneRequired       = errors.New("phone is required")
	ErrPhoneTooLong        = errors.New("phone max length 20")
	ErrApartmentRequired   = errors.New("apartment is required")
	ErrApartmentTooLong    = errors.New("apartment max length 20")
	ErrHouseIDRequired     = errors.New("house_id is required")
	ErrHouseIDInvalid      = errors.New("house_id must be a valid uuid")
	ErrCompanyNameRequired = errors.New("company name is required")
	ErrCompanyNameTooLong  = errors.New("company name max length 200")
	ErrHouseNameRequired   = errors.New("house name is required")
	ErrHouseAddressRequired = errors.New("house address is required")
	ErrUKIDRequired        = errors.New("uk_id is required")
)
