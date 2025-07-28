package main

import (
	"log"
	"net/http"

	"gridgod/backend"
)

func main() {
	log.Println("Backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", backend.NewHandler()))
}
