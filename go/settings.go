package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tarm/serial"
)

type Settings struct {
	Automatic bool    `json:"automatic"`
	Threshold float32 `json:"threshold"`
}

func HomeSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	buf, err := json.Marshal(settings)
	if err != nil {
		msg := fmt.Sprintf("failed to marshal json: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "failed to get Home settings state: %s"}`, msg)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func SetHomeSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	r.Body.Close()

	if err := json.Unmarshal(b, &settings); err != nil {
		msg := fmt.Sprintf("failed to unmarshal settings: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "failed to set Home settings state: %s"}`, msg)))
		return
	}

	if LIVE == true {
		s, err := serial.OpenPort(serialConf)
		if err != nil {
			msg := fmt.Sprintf("failed to open port: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Light toggle failed: %s"}\n`, msg)))
			return
		}

		// TODO: Should this be true/false?
		// Example Arduino commands: house_auto_ON, house_threshold_200
		var auto string
		if settings.Automatic {
			auto = "house_auto_ON\n"
		} else {
			auto = "house_auto_OFF\n"
		}

		log.Print("Sending:", auto)

		_, err = s.Write([]byte(auto))
		if err != nil {
			msg := fmt.Sprintf("failed to write: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Light toggle failed: %s"}\n`, msg)))
			return
		}

		threshold := fmt.Sprintf("house_threshold_%f\n", settings.Threshold)

		log.Print("Sending:", threshold)

		_, err = s.Write([]byte(threshold))
		if err != nil {
			msg := fmt.Sprintf("failed to write: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Light toggle failed: %s"}\n`, msg)))
			return
		}

		s.Close()
	}

	msg := fmt.Sprintf("OK, current house settings: automatic: '%t', threshold: '%f'",
		settings.Automatic, settings.Threshold)
	buf, err := json.Marshal(&StatusResponse{Message: msg})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}
