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
				Latitude:  float32Ptr(35.8436867),
				Longitude: float32Ptr(-78.7851406),
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
				Latitude:  float32Ptr(35.4757665),
				Longitude: float32Ptr(-80.79953),
			},
		},
	}
)

func float32Ptr(value float32) *float32 {
	return &value
}

// GetAllPropertiesHandler returns all properties
func GetAllPropertiesHandler(w http.ResponseWriter, r *http.Request) {
	client, ok := mongo.ClientFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to connect to db"))
	}
	db := client.Database("preservationnc")
	propsCollection := db.Collection("properties")
	cursor, err := propsCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not find any properties: %s", err.Error())))
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
		log.Printf("ERROR: %s", err.Error())
	}
	log.Printf("found properties = %v", props)
	err = mongo.SaveProperties(r.Context(), props)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to save properties: %s", err.Error())))
		return
	}
	w.WriteHeader(http.StatusCreated)
}
