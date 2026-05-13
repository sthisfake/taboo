package main

import (
	"log"
	"net/http"
	"os"
	"time"

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

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 15 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("github relay listening on :%s", port)

	log.Fatal(server.ListenAndServe())
}