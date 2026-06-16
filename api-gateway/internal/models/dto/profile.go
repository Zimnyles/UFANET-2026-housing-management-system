package dto

type ProfileResponse struct {
	UserID    string `json:"user_id"`
	FullName  string `json:"full_name"`
	Phone     string `json:"phone"`
	Apartment string `json:"apartment"`
	HouseID   string `json:"house_id"`
	UKID      string `json:"uk_id"`
}

type UpsertProfileRequest struct {
	FullName  string `json:"full_name"`
	Phone     string `json:"phone"`
	Apartment string `json:"apartment"`
	HouseID   string `json:"house_id"`
}

type ManagementCompanyResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateManagementCompanyRequest struct {
	Name string `json:"name"`
}

type ListManagementCompaniesResponse struct {
	Companies []ManagementCompanyResponse `json:"companies"`
}

type HouseResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	UKID    string `json:"uk_id"`
}

type CreateHouseRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	UKID    string `json:"uk_id"`
}

type ListHousesResponse struct {
	Houses []HouseResponse `json:"houses"`
}
