package domain

type Profile struct {
	UserID    string
	FullName  string
	Phone     string
	Apartment string
	HouseID   string
	UKID      string
}

type UpsertProfileRequest struct {
	UserID    string
	FullName  string
	Phone     string
	Apartment string
	HouseID   string
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

type CreateManagementCompanyRequest struct {
	Name string
}

type CreateHouseRequest struct {
	Name    string
	Address string
	UKID    string
}
