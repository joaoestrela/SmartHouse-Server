package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
		req := "luminosity\n"
		log.Print("sending:", req)

		mutex.Lock()

		_, err = s.Write([]byte(req))
		if err != nil {
			msg := fmt.Sprintf("failed to write: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Luminosity get failed: %s"}`, msg)))
			return
		}

		// Read response with newline delim
		reader := bufio.NewReader(s)
		val, err := reader.ReadBytes('\n')

		mutex.Unlock()

		if err != nil {
			msg := fmt.Sprintf("failed to read: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Luminosity get failed: %s"}`, msg)))
			return
		}

		parsed := strings.TrimSpace(string(val))
		log.Println("incoming:", parsed)

		lum, err := strconv.ParseFloat(parsed, 32)
		if err != nil {
			msg := fmt.Sprintf("failed to parse: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Luminosity get failed: %s"}`, msg)))
			return
		}

		sd.Value = float32(lum)
		sd.Unit = "Lux"

	} else {
		sd = SensorData{
			Value: 1.2,
			Unit:  "Lux",
		}
	}

	log.Printf("Value: %f, Unit: %s\n", sd.Value, sd.Unit)

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
		req := "temperature\n"
		log.Print("sending:", req)

		mutex.Lock()

		_, err = s.Write([]byte(req))
		if err != nil {
			msg := fmt.Sprintf("failed to write: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Temperature get failed: %s"}`, msg)))
			return
		}

		// Try to read response
		log.Println("reading from serial")

		reader := bufio.NewReader(s)
		val, err := reader.ReadBytes('\n')

		mutex.Unlock()

		if err != nil {
			msg := fmt.Sprintf("failed to read: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Temperature get failed: %s"}`, msg)))
			return
		}

		parsed := strings.TrimSpace(string(val))
		log.Println("incoming:", parsed)

		temp, err := strconv.ParseFloat(parsed, 32)
		if err != nil {
			msg := fmt.Sprintf("failed to parse: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Temperature get failed: %s"}`, msg)))
			return
		}

		sd.Value = float32(temp)
		sd.Unit = "Celcius"

	} else {
		sd = SensorData{
			Value: 20,
			Unit:  "Celcius",
		}
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
