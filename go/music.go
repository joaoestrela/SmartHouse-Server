package server

import (
	"encoding/json"
	"net/http"
)

type MusicPlayerStatus struct {
	State bool `json:"state,omitempty"`
	Track Track
}

type Track struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func MusicAvailable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	tracks := []Track{
		Track{
			ID:   42,
			Name: "Rick Astley - Never Gonna Give You Up",
		},
		Track{
			ID:   43,
			Name: "Air Supply - All Out of Love",
		},
	}
	buf, err := json.Marshal(tracks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)

}

func MusicSummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	status := MusicPlayerStatus{
		State: true,
		Track: Track{
			ID:   42,
			Name: "Rick Astley - Never Gonna Give You Up",
		},
	}
	buf, err := json.Marshal(status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func PlayTrack(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	buf, err := json.Marshal(StatusResponse{Message: "OK"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func SetMusicState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	status := MusicPlayerStatus{
		State: true,
		Track: Track{
			ID:   42,
			Name: "Rick Astley - Never Gonna Give You Up",
		},
	}
	buf, err := json.Marshal(status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}
