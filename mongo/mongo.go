package mongo

import (
	"context"
	"fmt"
	"log"

	"github.com/toast2e/preservationnc-server/reps"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type contextKey string

const (
	clientContextKey contextKey = contextKey("mongodb:client")
)

// NewClientContext creates a new context with a mongo client
func NewClientContext(ctx context.Context, client *mongo.Client) context.Context {
	return context.WithValue(ctx, clientContextKey, client)
}

// ClientFromContext gets the mongo client from the current context if available
func ClientFromContext(ctx context.Context) (*mongo.Client, bool) {
	client, ok := ctx.Value(clientContextKey).(*mongo.Client)
	return client, ok
}

// SetupClient sets up and connects a mongo client and stores it in a new context derived from the given context
func SetupClient(ctx context.Context) (context.Context, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://db:27017"))
	if err != nil {
		return ctx, err
	}
	err = client.Connect(ctx)
	if err != nil {
		return ctx, err
	}

	return NewClientContext(ctx, client), nil
}

// Shutdown disconnects the mongo client found in the context if available
func Shutdown(ctx context.Context) {
	client, ok := ClientFromContext(ctx)
	if ok {
		client.Disconnect(ctx)
	}
}

// SaveProperties saves the properties given using the client found in the given context
func SaveProperties(ctx context.Context, properties []reps.Property) error {
	client, ok := ClientFromContext(ctx)
	if !ok {
		return fmt.Errorf("no mongo client found in context: %v", ctx)
	}
	db := client.Database("preservationnc")
	propsCollection := db.Collection("properties")
	for _, prop := range properties {
		res, err := propsCollection.InsertOne(ctx, prop)
		if err != nil {
			log.Printf("WARN: failed to store property %s (%s): %s", prop.ID, prop.Name, err.Error())
		} else {
			log.Printf("DB_ID for property ID = %s: %s", prop.ID, res.InsertedID)
		}
	}
	return nil
}
