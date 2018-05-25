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

	buf, err := json.Marshal(tracks)
	if err != nil {
		msg := fmt.Sprintf("failed to marshal json: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "Available failed: %s"}`, msg)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)

}

func MusicSummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	status := MusicPlayerStatus{
		State: trackPlaying,
		Track: activeTrack,
	}
	buf, err := json.Marshal(status)
	if err != nil {
		msg := fmt.Sprintf("failed to marshal json: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "Summary failed: %s"}`, msg)))
		return
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

	trackPlaying = true
	activeTrack = tracks[t.ID-1]

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message": "OK, playing track: %d - %s"}`,
		activeTrack.ID, activeTrack.Name)))
}

func SetMusicState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	params := mux.Vars(r)
	state := params["state"]

	if state == "on" {
		trackPlaying = true
	} else if state == "off" {
		trackPlaying = false
	} else {
		msg := fmt.Sprintf("unknown music state: %v", state)
		log.Println(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message": "Music state failed: %s"}`, msg)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message": "OK, music state updated to: %s"}`, state)))
}
