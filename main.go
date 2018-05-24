package main

import (
	"log"
	"net/http"

	server "github.com/freddygv/SmartHouse-Server/go"
)

var (
	device = "/dev/ttyACM0"
	baud   = 9600
)

func main() {
	router := server.NewServer(device, baud)

	log.Printf("Server started...\n")
	log.Fatal(http.ListenAndServe("0.0.0.0:8888", router))

	// Make requests for temperature
	// Make requests for luminosity
	// Make request for settings
}
