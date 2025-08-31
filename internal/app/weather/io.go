package weather

type FetchDataInput struct {
	CityName string `json:"city_name" validate:"required"`
	Country  string `json:"country"`
}
