package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/gorilla/mux"
)

type MusicPlayerStatus struct {
	State bool
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

	id, err := strconv.Atoi(r.URL.Query().Get("trackId"))
	if err != nil {
		msg := fmt.Sprintf("failed to parse id: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "Play failed: %s"}`, msg)))
		return
	}

	if LIVE {
		mpg123 = exec.Command("mpg123", "-q", music+tracks[id-1].Name)
		err = mpg123.Start()
		if err != nil {
			msg := fmt.Sprintf("failed to play: %v", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Play failed: %s"}`, msg)))
			return
		}
	}

	trackPlaying = true
	activeTrack = tracks[id-1]

	status := MusicPlayerStatus{
		State: trackPlaying,
		Track: activeTrack,
	}
	buf, err := json.Marshal(status)
	if err != nil {
		msg := fmt.Sprintf("failed to marshal json: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "Play failed: %s"}`, msg)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func SetMusicState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	params := mux.Vars(r)
	state := params["state"]

	if state == "off" {
		trackPlaying = false
		activeTrack = Track{}

		if LIVE && mpg123 != nil {
			mpg123.Process.Kill()
			mpg123 = nil
		}

	} else {
		msg := fmt.Sprintf("unknown music state: %v", state)
		log.Println(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message": "Music state failed: %s"}`, msg)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "OK, music state updated to off"}`))
}
