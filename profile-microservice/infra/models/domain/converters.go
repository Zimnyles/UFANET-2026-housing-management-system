package domain

import "profile-service/infra/models/dto"

func UpsertProfileRequestFromDTO(d *dto.UpsertProfileRequest) *UpsertProfileRequest {
	return &UpsertProfileRequest{
		UserID:    d.UserID,
		FullName:  d.FullName,
		Phone:     d.Phone,
		Apartment: d.Apartment,
		HouseID:   d.HouseID,
	}
}

func ProfileToDTO(p *Profile) *dto.Profile {
	return &dto.Profile{
		UserID:    p.UserID,
		FullName:  p.FullName,
		Phone:     p.Phone,
		Apartment: p.Apartment,
		HouseID:   p.HouseID,
		UKID:      p.UKID,
	}
}

func ManagementCompanyToDTO(c *ManagementCompany) *dto.ManagementCompany {
	return &dto.ManagementCompany{ID: c.ID, Name: c.Name}
}

func HouseToDTO(h *House) *dto.House {
	return &dto.House{ID: h.ID, Name: h.Name, Address: h.Address, UKID: h.UKID}
}
