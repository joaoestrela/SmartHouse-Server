package server

import (
	"encoding/json"
	"net/http"
)

type Settings struct {
	Automatic bool    `json:"automatic,omitempty"`
	Threshold float32 `json:"threshold,omitempty"`
}

func SetHomeSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	buf, err := json.Marshal(&StatusResponse{Message: "OK"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}
