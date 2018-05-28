package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/gorilla/mux"
	"github.com/tarm/serial"
)

const (
	LIVE    = true
	storage = "auth.db"
	music   = "/home/pi/music/"
)

var (
	serialConf   *serial.Config
	lights       []Light
	tracks       []Track
	activeTrack  Track
	trackPlaying bool
	mpg123       *exec.Cmd
	temperature  float32
	luminosity   float32
	settings     Settings
)

type StatusResponse struct {
	Message string `json:"message,omitempty"`
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewServer(device string, baud int) *mux.Router {
	if err := NewAuthDB(storage); err != nil {
		log.Fatalf("failed to create db: %v", err)
	}

	// Configure serial comms
	serialConf = &serial.Config{Name: device, Baud: baud}

	// Init Light state for each bedroom, all light start off
	rooms := []string{"bedroom-1", "bedroom-2", "living room", "kitchen", "bathroom"}

	for i := 0; i < 5; i++ {
		l := Light{
			ID:          i + 1,
			Description: rooms[i],
			TurnOn:      false,
		}

		lights = append(lights, l)
	}

	// Init Settings
	settings = Settings{
		Automatic: false,
		Threshold: 1,
	}

	if LIVE {
		// Init songs
		files, err := ioutil.ReadDir(music)
		if err != nil {
			log.Fatalf("failed to read music dir: %v", err)
		}

		for i, f := range files {
			t := Track{
				ID:   i + 1,
				Name: f.Name(),
			}
			tracks = append(tracks, t)
		}

		// Launch receiver for serial updates from Arduino
		go UpdateReceiver()

	} else {
		tracks = []Track{
			Track{
				ID:   1,
				Name: "Rick Astley - Never Gonna Give You Up",
			},
			Track{
				ID:   2,
				Name: "Rick Astley - Whenever You Need Somebody",
			},
			Track{
				ID:   3,
				Name: "Rick Astley - Together Forever",
			},
			Track{
				ID:   4,
				Name: "Rick Astley - It Would Take a Strong Strong Man",
			},
			Track{
				ID:   5,
				Name: "Rick Astley - The Love Has Gone",
			},
			Track{
				ID:   6,
				Name: "Rick Astley - Don't Say Goodbye",
			},
			Track{
				ID:   7,
				Name: "Rick Astley - Slipping Away",
			},
			Track{
				ID:   8,
				Name: "Rick Astley - No More Looking for Love",
			},
			Track{
				ID:   9,
				Name: "Rick Astley - You Move Me",
			},
			Track{
				ID:   10,
				Name: "Rick Astley - When I Fall in Love",
			},
		}
	}

	// Init routes
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}

	return router
}

func Health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

var routes = Routes{
	Route{
		"Health",
		"GET",
		"/SmartHouse/1.0.2/health",
		Health,
	},

	Route{
		"Login",
		"POST",
		"/SmartHouse/1.0.2/login",
		Login,
	},

	Route{
		"Register",
		"POST",
		"/SmartHouse/1.0.2/register",
		Register,
	},

	Route{
		"Luminosity",
		"GET",
		"/SmartHouse/1.0.2/luminosity",
		Luminosity,
	},

	Route{
		"Temperature",
		"GET",
		"/SmartHouse/1.0.2/temperature",
		Temperature,
	},

	Route{
		"LightState",
		"GET",
		"/SmartHouse/1.0.2/lights/{lightID}",
		LightState,
	},

	Route{
		"Lights",
		"GET",
		"/SmartHouse/1.0.2/lights",
		Lights,
	},

	Route{
		"SetLightState",
		"PUT",
		"/SmartHouse/1.0.2/lights/{lightID}/{state}",
		SetLightState,
	},

	Route{
		"MusicAvailable",
		"GET",
		"/SmartHouse/1.0.2/music/available/",
		MusicAvailable,
	},

	Route{
		"MusicSummary",
		"GET",
		"/SmartHouse/1.0.2/music",
		MusicSummary,
	},

	Route{
		"PlayTrack",
		"PUT",
		"/SmartHouse/1.0.2/music/play",
		PlayTrack,
	},

	Route{
		"SetMusicState",
		"PUT",
		"/SmartHouse/1.0.2/music/{state}",
		SetMusicState,
	},

	Route{
		"HomeSettings",
		"GET",
		"/SmartHouse/1.0.2/settings/home/",
		HomeSettings,
	},

	Route{
		"SetHomeSettings",
		"PUT",
		"/SmartHouse/1.0.2/settings/home/",
		SetHomeSettings,
	},
}

func UpdateReceiver() {
	for {
		s, err := serial.OpenPort(serialConf)
		if err != nil {
			log.Fatalf("failed to open serial port: %v", err)
		}

		// Try to read
		dec := json.NewDecoder(s)
		for dec.More() {
			var l Light
			if err := dec.Decode(&l); err != nil {
				log.Printf("error: failed to unmarshal: %v", err)
				time.Sleep(1 * time.Millisecond)
				continue
			}

			lights[l.ID-1].TurnOn = l.TurnOn
			log.Printf("Light #%d set to: %t", l.ID, l.TurnOn)
		}
		s.Close()
	}
}
