package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tarm/serial"
)

type SensorData struct {
	Value float32 `json:"value"`
	Unit  string  `json:"unit,omitempty"`
}

func Luminosity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var sd SensorData

	if LIVE == true {
		// Open port for 2-way comms
		s, err := serial.OpenPort(serialConf)
		if err != nil {
			msg := fmt.Sprintf("failed to open port: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Luminosity get failed: %s"}`, msg)))
			return
		}
		defer s.Close()

		// Write request
		_, err = s.Write([]byte("luminosity"))
		if err != nil {
			msg := fmt.Sprintf("failed to write: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Luminosity get failed: %s"}`, msg)))
			return
		}

		// Read response with newline delim
		reader := bufio.NewReader(s)
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			msg := fmt.Sprintf("failed to read: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Luminosity get failed: %s"}`, msg)))
			return
		}

		log.Printf("incoming: %s\n", string(msg))

		// Unmarshal response into struct
		if err := json.Unmarshal(msg, &sd); err != nil {
			msg := fmt.Sprintf("failed to unmarshal: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Luminosity get failed: %s"}`, msg)))
			return
		}

		// TODO: Verify
		sd.Unit = "Lux"

	} else {
		sd = SensorData{
			Value: 1.2,
			Unit:  "Lux",
		}
	}

	buf, err := json.Marshal(sd)
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
			Unit:  "Lux",
		},
		SensorData{
			Value: 1.0,
			Unit:  "Lux",
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
	var sd SensorData

	if LIVE == true {
		s, err := serial.OpenPort(serialConf)
		if err != nil {
			msg := fmt.Sprintf("failed to open port: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Temperature get failed: %s"}`, msg)))
			return
		}
		defer s.Close()

		// Write request
		_, err = s.Write([]byte("temperature"))
		if err != nil {
			msg := fmt.Sprintf("failed to write: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Temperature get failed: %s"}`, msg)))
			return
		}

		// Try to read response
		reader := bufio.NewReader(s)
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			msg := fmt.Sprintf("failed to read: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Temperature get failed: %s"}`, msg)))
			return
		}
		log.Printf("incoming: %s\n", string(msg))

		if err := json.Unmarshal(msg, &sd); err != nil {
			msg := fmt.Sprintf("failed to read: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Temperature get failed: %s"}`, msg)))
			return
		}

		// TODO: Verify
		sd.Unit = "Celcius"

	} else {
		sd = SensorData{
			Value: 20,
			Unit:  "Celcius",
		}
	}

	buf, err := json.Marshal(sd)
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
