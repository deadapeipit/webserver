package entity

type Address struct {
	City          string      `json:"city"`
	StreetName    string      `json:"street_name"`
	StreetAddress string      `json:"street_address"`
	ZipCode       string      `json:"zip_code"`
	State         string      `json:"state"`
	Country       string      `json:"country"`
	Coordinates   Coordinates `json:"coordinates"`
}
