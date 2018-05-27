package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SensorData struct {
	Value float32 `json:"value"`
	Unit  string  `json:"unit,omitempty"`
}

func Luminosity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	sd := SensorData{
		Value: 688,
		Unit:  "Lux",
	}

	buf, err := json.Marshal(sd)
	if err != nil {
		msg := fmt.Sprintf("failed to marshal: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "Luminosity get failed: %s"}`, msg)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func Temperature(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	sd := SensorData{
		Value: 27,
		Unit:  "Celcius",
	}

	buf, err := json.Marshal(sd)
	if err != nil {
		msg := fmt.Sprintf("failed to marshal: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "Temperature get failed: %s"}`, msg)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}
