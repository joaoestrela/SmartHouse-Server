package main

import (
	"log"
	"net/http"

	server "github.com/freddygv/SmartHouse-Server/go"
)

func main() {
	router := server.NewRouter()

	log.Printf("Server started")
	log.Fatal(http.ListenAndServe("0.0.0.0:8888", router))
}
