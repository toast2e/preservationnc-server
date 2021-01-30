package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/toast2e/preservationnc-server/html"
	"github.com/toast2e/preservationnc-server/mongo"
	"github.com/toast2e/preservationnc-server/reps"
	"go.mongodb.org/mongo-driver/bson"
	"googlemaps.github.io/maps"
)

var (
	// DummyProps are used for testing purposes
	DummyProps = []reps.Property{
		{
			Name:        "raleighProperty1",
			Description: "testDescription",
			Price:       200000.00,
			Location: reps.Site{
				Address:   "123 Fake Street",
				City:      "Raleigh",
				State:     "North Carolina",
				Zip:       "12345",
				Latitude:  float64Ptr(35.8436867),
				Longitude: float64Ptr(-78.7851406),
			},
		},
		{
			Name:        "kannapolisProperty1",
			Description: "this property is in kannapolis",
			Price:       100000.00,
			Location: reps.Site{
				Address:   "321 teertS ekaF",
				City:      "Kannpolis",
				State:     "North Carolina",
				Zip:       "54321",
				Latitude:  float64Ptr(35.4757665),
				Longitude: float64Ptr(-80.79953),
			},
		},
	}
)

func float64Ptr(value float64) *float64 {
	return &value
}

// GetAllPropertiesHandler returns all properties
func GetAllPropertiesHandler(w http.ResponseWriter, r *http.Request) {
	client, ok := mongo.ClientFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to connect to db"))
		return
	}
	db := client.Database("preservationnc")
	propsCollection := db.Collection("properties")
	cursor, err := propsCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not find any properties: %s", err.Error())))
		return
	}

	// TODO size this properly
	properties := make([]reps.Property, 0)
	err = cursor.All(context.TODO(), &properties)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to unmarshal properties from bson: %s", err.Error())))
		return
	}
	propertiesJSON, err := json.Marshal(properties)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to marshal properties to json: %s", err.Error())))
		return
	}

	w.Write(propertiesJSON)
}

// DeleteAll deletes all properties
func DeleteAll(w http.ResponseWriter, r *http.Request) {
	client, ok := mongo.ClientFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to connect to db"))
		return
	}
	db := client.Database("preservationnc")
	propsCollection := db.Collection("properties")
	res, err := propsCollection.DeleteMany(r.Context(), bson.D{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to delete all properties: %s", err.Error())))
		return
	}
	w.Write([]byte(fmt.Sprintf("deleted %d properties", res.DeletedCount)))
}

// Reload reloads all properties from the source
func Reload(w http.ResponseWriter, r *http.Request) {
	client := http.Client{Timeout: 10 * time.Second}
	crawler := html.NewCrawler(client)
	props, err := crawler.FindProperties()
	if err != nil {
		log.Printf("ERROR - fetching properties: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to fetch properties: %s", err.Error())))
		//return
	}
	// TODO configure a single client for the app at startup and add to the context
	mapsClient, err := maps.NewClient(maps.WithAPIKey("<insert API key here>"))
	if err != nil {
		log.Printf("ERROR - configuring maps client: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed connecting to maps API: %s", err.Error())))
		return
	}

	for i, prop := range props {
		geoRequest := maps.GeocodingRequest{Address: fmt.Sprintf("%s, %s, %s, %s", prop.Location.Address, prop.Location.City, prop.Location.State, prop.Location.Zip)}
		results, err := mapsClient.Geocode(r.Context(), &geoRequest)
		if err != nil {
			log.Printf("ERROR - retrieving property geo info: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("failed retrieving property geo info: %s", err.Error())))
			return
		}
		log.Printf("got geo info: %v", results[0])
		r := results[0]
		prop.Location.Latitude = &r.Geometry.Location.Lat
		prop.Location.Longitude = &r.Geometry.Location.Lng
		log.Printf("%s: %d, %d", prop.Name, prop.Location.Latitude, prop.Location.Longitude)
		props[i] = prop
	}

	log.Printf("found properties = %v", props)
	err = mongo.SaveProperties(r.Context(), props)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("ERROR - failed to save properties: %s", err.Error())))
		return
	}
	w.WriteHeader(http.StatusCreated)
}
