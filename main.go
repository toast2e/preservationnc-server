package main

import (
	"log"
	"net/http"

	phttp "github.com/toast2e/preservationnc-server/http"
)

func main() {
	http.HandleFunc("/properties", phttp.GetAllPropertiesHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
