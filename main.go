package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	phttp "github.com/toast2e/preservationnc-server/http"
	pmongo "github.com/toast2e/preservationnc-server/mongo"
)

func main() {
	ctx := context.Background()
	ctx, err := pmongo.SetupClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup db client: %s", err.Error())
	}
	defer pmongo.Shutdown(ctx)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	route := os.Getenv("BASE_PATH")
	if route == "" {
		route = "/preservationnc"
	}

	// TODO uncomment after figuring out db stuff
	//client := http.Client{Timeout: 10 * time.Second}
	//crawler := phtml.NewCrawler(client)
	//props, err := crawler.FindProperties()
	//if err != nil {
	//	log.Printf("ERROR: %s", err.Error())
	//}
	//log.Printf("found properties = %v", props)
	//pmongo.SaveProperties(ctx, props)

	// TODO remove this after figuring out db stuff
	err = pmongo.SaveProperties(ctx, phttp.DummyProps)
	if err != nil {
		log.Fatalf("Could not load initial properties: %s", err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("%s/properties", route), phttp.GetAllPropertiesHandler)
	mux.HandleFunc(fmt.Sprintf("%s/delete", route), phttp.DeleteAll)
	mux.HandleFunc(fmt.Sprintf("%s/reload", route), phttp.Reload)
	updatedMux := addClientToRequestContext(ctx, mux)

	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, updatedMux))
}

func addClientToRequestContext(ctx context.Context, next http.Handler) http.Handler {
	client, _ := pmongo.ClientFromContext(ctx)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(pmongo.NewClientContext(r.Context(), client))
		next.ServeHTTP(w, r)
	})
}
