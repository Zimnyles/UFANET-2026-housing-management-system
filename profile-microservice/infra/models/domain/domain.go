package domain

import "time"

type Profile struct {
	UserID    string
	FullName  string
	Phone     string
	Apartment string
	HouseID   string
	UKID      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ManagementCompany struct {
	ID   string
	Name string
}

type House struct {
	ID      string
	Name    string
	Address string
	UKID    string
}

type UpsertProfileRequest struct {
	UserID    string
	FullName  string
	Phone     string
	Apartment string
	HouseID   string
}

type CreateManagementCompanyRequest struct {
	Name string
}

type CreateHouseRequest struct {
	Name    string
	Address string
	UKID    string
}

type ListHousesRequest struct {
	UKID string
}
