package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	handler "go-xhttp-relay/api"
)

func main() {
	_ = godotenv.Load()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
