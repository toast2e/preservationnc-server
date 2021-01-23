package reps

// Property contains all values related to an individual property
type Property struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Location    Site    `json:"location"`
}

// Site represents where a property is located
type Site struct {
	Address   string   `json:"address"`
	City      string   `json:"city"`
	County    string   `json:"county"`
	State     string   `json:"state"`
	Zip       string   `json:"zip"`
	Longitude *float32 `json:"longitude,omitempty"`
	Latitude  *float32 `json:"latitude,omitempty"`
}
