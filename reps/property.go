package reps

type Property struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Location    Site    `json:"location"`
}

type Site struct {
	Address   Address `json:"address"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Zip       string  `json:"zip"`
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

type Address struct {
	Number          string `json:"number"`
	ApartmentNumber string `json:"apartmentNumber"`
	Street          string `json:"street"`
}
