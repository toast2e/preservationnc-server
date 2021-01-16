package http

import (
	"encoding/json"
	"net/http"

	"github.com/toast2e/preservationnc-server/reps"
)

var (
	dummyProps = []reps.Property{
		{
			Name:        "raleighProperty1",
			Description: "testDescription",
			Price:       200000.00,
			Location: reps.Site{
				// Address: reps.Address{
				// 	Number: "123",
				// 	Street: "Fake Street",
				// },
				Address:   "123 Fake Street",
				City:      "Raleigh",
				State:     "North Carolina",
				Zip:       "12345",
				Latitude:  35.8436867,
				Longitude: -78.7851406,
			},
		},
		{
			Name:        "kannapolisProperty1",
			Description: "this property is in kannapolis",
			Price:       100000.00,
			Location: reps.Site{
				// Address: reps.Address{
				// 	Number: "321",
				// 	Street: "teertS ekaF",
				// },
				Address:   "321 teertS ekaF",
				City:      "Kannpolis",
				State:     "North Carolina",
				Zip:       "54321",
				Latitude:  35.4757665,
				Longitude: -80.79953,
			},
		},
	}
)

func GetAllPropertiesHandler(w http.ResponseWriter, r *http.Request) {
	propertiesJson, err := json.Marshal(dummyProps)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to marshal json"))
	}
	w.Write(propertiesJson)
}
