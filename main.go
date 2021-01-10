package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	phttp "github.com/toast2e/preservationnc-server/http"
	"github.com/toast2e/preservationnc-server/reps"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	ctx = setupDB(ctx)
	defer shutdown(ctx)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	route := os.Getenv("BASE_PATH")
	if route == "" {
		route = "/preservationnc"
	}

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://www.presnc.org/property-listing/all-properties/")
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal(fmt.Sprintf("unexpected response from server: %s", resp.Status))
	}
	bodyBytes := make([]byte, 1000000)
	n, err := resp.Body.Read(bodyBytes)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("read %d bytes from server\n", n)
	log.Printf("%s", string(bodyBytes))

	http.HandleFunc(fmt.Sprintf("%s/properties", route), phttp.GetAllPropertiesHandler)
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func setupDB(ctx context.Context) context.Context {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://db:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("preservationnc")
	properties := db.Collection("properties")
	res, err := properties.InsertOne(ctx, reps.Property{
		Name:        "raleighProperty1",
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
			Longitude: -78.7851406,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %s\n", res.InsertedID)

	return context.WithValue(ctx, "mongodb:client", client)
}

func shutdown(ctx context.Context) {
	mongoClient, ok := ctx.Value("mongodb:client").(mongo.Client)
	if ok {
		mongoClient.Disconnect(ctx)
	}
}
