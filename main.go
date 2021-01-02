package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	phttp "github.com/toast2e/preservationnc-server/http"
)

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	route := os.Getenv("BASE_PATH")
	if route == "" {
		route = "/preservationnc"
	}

	http.HandleFunc(fmt.Sprintf("%s/properties", route), phttp.GetAllPropertiesHandler)
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
