package server

import (
	"net/http"
)

type Settings struct {
	Automatic bool    `json:"automatic,omitempty"`
	Threshold float32 `json:"threshold,omitempty"`
}

func SetHomeSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
