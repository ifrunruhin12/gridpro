package main

import (
	"log"
	"net/http"
	"os"

	"gridpro/backend"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, backend.NewHandler()))
}

