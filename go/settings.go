package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Settings struct {
	Automatic bool    `json:"automatic,omitempty"`
	Threshold float32 `json:"threshold,omitempty"`
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
	defer r.Body.Close()

	if err := json.Unmarshal(b, &settings); err != nil {
		msg := fmt.Sprintf("failed to unmarshal settings: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "failed to set Home settings state: %s"}`, msg)))
		return
	}

	// TODO: Write two commands to Arduino, one for threshold and one for automatic

	msg := fmt.Sprintf("OK, current house settings: automatic: '%t', threshold: '%f'",
		settings.Automatic, settings.Threshold)
	buf, err := json.Marshal(&StatusResponse{Message: msg})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}
