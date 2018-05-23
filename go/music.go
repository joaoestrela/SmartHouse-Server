package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()

	var t Track
	if err := json.Unmarshal(b, &t); err != nil {
		msg := fmt.Sprintf("failed to unmarshal: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "Play failed: %s"}`, msg)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message": "OK, playing track id: %s"}`, t.ID)))
}

func SetMusicState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	params := mux.Vars(r)
	state := params["state"]

	// TODO: Do something with request

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message": "OK, music state updated to: %s"}`, state)))
}
