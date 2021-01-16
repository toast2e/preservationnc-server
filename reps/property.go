package reps

type Property struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Location    Site    `json:"location"`
}

type Site struct {
	Address   string  `json:"address"`
	City      string  `json:"city"`
	County    string  `json:"county"`
	State     string  `json:"state"`
	Zip       string  `json:"zip"`
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

// type Address struct {
// 	Number          string `json:"number"`
// 	ApartmentNumber string `json:"apartmentNumber"`
// 	Street          string `json:"street"`
// }
