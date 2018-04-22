package server

import (
	"encoding/json"
	"net/http"
)

type Light struct {
	ID          int     `json:"id,omitempty"`
	Description string  `json:"description,omitempty"`
	On          bool    `json:"on,omitempty"`
	Threshold   float32 `json:"threshold,omitempty"`
	Automatic   bool    `json:"automatic,omitempty"`
}

func GetLightState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	buf, err := json.Marshal(&Light{
		ID:          1,
		Description: "bedroom",
		On:          false,
		Threshold:   0.5,
		Automatic:   true,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func GetLights(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	lights := []Light{
		Light{
			ID:          1,
			Description: "bedroom",
			On:          false,
			Threshold:   0.5,
			Automatic:   true,
		},
		Light{
			ID:          1,
			Description: "kitchen",
			On:          true,
			Threshold:   0.5,
			Automatic:   true,
		},
	}
	buf, err := json.Marshal(lights)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func SetLightState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	buf, err := json.Marshal(&StatusResponse{Message: "OK"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func SetLuminosityThreshhold(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	buf, err := json.Marshal(&StatusResponse{Message: "OK"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func SettingsLight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	buf, err := json.Marshal(&Settings{Automatic: true, Threshold: 1.5})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}
