package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/tarm/serial"
)

type Light struct {
	ID          int    `json:"id"`
	Description string `json:"description,omitempty"`
	TurnOn      bool   `json:"turnon"`
}

func LightState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["lightID"])
	if err != nil {
		msg := fmt.Sprintf("failed to parse id: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "Light state failed: %s"}`, msg)))
		return
	}

	buf, err := json.Marshal(lights[id])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func Lights(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// TODO: Verify what's going on here
	fmt.Println(lights[2].TurnOn)

	buf, err := json.Marshal(lights)
	if err != nil {
		msg := fmt.Sprintf("failed to marshal json: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "Lights state failed: %s"}`, msg)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func SetLightState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	params := mux.Vars(r)

	state := strings.ToUpper(params["state"])

	id, err := strconv.Atoi(params["lightID"])
	if err != nil || id >= len(lights)-1 {
		msg := fmt.Sprintf("failed to parse id: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "Light toggle failed: %s"}\n`, msg)))
		return
	}

	switch state {
	case "ON":
		lights[id].TurnOn = true
	case "OFF":
		lights[id].TurnOn = false
	default:
		msg := fmt.Sprintf("invalid command: %v", state)
		log.Println(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message": "Light toggle failed: %s"}`, msg)))
		return
	}

	if LIVE == true {
		s, err := serial.OpenPort(serialConf)
		if err != nil {
			msg := fmt.Sprintf("failed to open port: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Light toggle failed: %s"}`, msg)))
			return
		}
		defer s.Close()

		// Example Arduino commands: led1_ON, led2_OFF
		cmd := fmt.Sprintf("led%d_%s\n", id, state)

		_, err = s.Write([]byte(cmd))
		if err != nil {
			msg := fmt.Sprintf("failed to write: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Light toggle failed: %s"}`, msg)))
			return
		}
	}

	msg := fmt.Sprintf("OK, toggled light #%d to %s", id, state)
	buf, err := json.Marshal(&StatusResponse{Message: msg})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}
