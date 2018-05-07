package main

import (
	"log"
	"net/http"

	server "github.com/freddygv/SmartHouse-Server/go"
)

const storage = "auth.db"

func main() {
	if err := server.NewAuthDB(storage); err != nil {
		log.Fatalf("failed to create db: %v", err)
	}
	router := server.NewRouter()

	log.Printf("Server started...\n")
	log.Fatal(http.ListenAndServe("0.0.0.0:8888", router))
}
