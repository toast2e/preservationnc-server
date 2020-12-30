package http

import (
	"encoding/json"
	"net/http"

	"github.com/toast2e/preservationnc-server/reps"
)

var (
	dummyProps = []reps.Property{
		{
			Name:        "testProperty1",
			Description: "testDescription",
			Price:       200000.00,
			Location: reps.Site{
				Address: reps.Address{
					Number: "123",
					Street: "Fake Street",
				},
				City:      "Raleigh",
				State:     "North Carolina",
				Zip:       "12345",
				Latitude:  35.8436867,
				Longitute: -78.7851406,
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
