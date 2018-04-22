package server

import (
	"encoding/json"
	"net/http"
)

type SensorData struct {
	Value float32 `json:"value,omitempty"`
	Unit  string  `json:"unit,omitempty"`
}

func Luminosity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	buf, err := json.Marshal(SensorData{
		Value: 1.2,
		Unit:  "candelas",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func LuminosityHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	hist := []SensorData{
		SensorData{
			Value: 1.2,
			Unit:  "candelas",
		},
		SensorData{
			Value: 1.0,
			Unit:  "candelas",
		},
	}
	buf, err := json.Marshal(hist)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func Temperature(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	buf, err := json.Marshal(SensorData{
		Value: 20,
		Unit:  "Celsius",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func TemperatureHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	hist := []SensorData{
		SensorData{
			Value: 20,
			Unit:  "Celsius",
		},
		SensorData{
			Value: 23,
			Unit:  "Celsius",
		},
	}
	buf, err := json.Marshal(hist)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}
