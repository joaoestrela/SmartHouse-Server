package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tarm/serial"
)

const (
	LIVE    = false
	storage = "auth.db"
)

var (
	serialConf *serial.Config
	lights     []Light
	settings   Settings
)

type StatusResponse struct {
	Message string `json:"message,omitempty"`
}

type Update struct {
	Source string `json:"source,omitempty"`
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
			ID:          i,
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

	// Launch receiver for serial updates from Arduino
	if LIVE {
		go UpdateReceiver()
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
		"LuminosityHistory",
		"GET",
		"/SmartHouse/1.0.2/luminosity/history",
		LuminosityHistory,
	},

	Route{
		"Temperature",
		"GET",
		"/SmartHouse/1.0.2/temperature",
		Temperature,
	},

	Route{
		"TemperatureHistory",
		"GET",
		"/SmartHouse/1.0.2/temperature/history",
		TemperatureHistory,
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
	s, err := serial.OpenPort(serialConf)
	if err != nil {
		log.Fatalf("failed to open serial port: %v", err)
	}
	defer s.Close()

	for {
		// Try to read
		reader := bufio.NewReader(s)
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("error: failed to read: %v\n", err)
			time.Sleep(1 * time.Second) // TODO: debug
			continue
		}
		log.Printf("incoming: %s\n", string(msg))

		var u Update
		if err := json.Unmarshal(msg, &u); err != nil {
			log.Printf("error: failed to unmarshal: %v", err)
			continue
		}

		switch u.Source {
		case "light":
			var l Light
			if err := json.Unmarshal(msg, &l); err != nil {
				log.Printf("error: failed to unmarshal: %v", err)
				continue
			}

			lights[l.ID].TurnOn = l.TurnOn
			log.Printf("Light #%d set to: %t", l.ID, l.TurnOn)

		case "settings":
			var s Settings
			if err := json.Unmarshal(msg, &s); err != nil {
				log.Printf("error: failed to unmarshal: %v", err)
				continue
			}

			// TODO: Import state into settings struct
			settings.Automatic = s.Automatic
			settings.Threshold = s.Threshold

			log.Printf("New settings: automatic: '%t', threshold: '%f'",
				settings.Automatic, settings.Threshold)

		case "sensor":
			var s SensorData
			if err := json.Unmarshal(msg, &s); err != nil {
				log.Printf("error: failed to unmarshal: %v", err)
				continue
			}

			// TODO: Save to DB
		}

		// TODO: Make sure arduino can receive
		// Reply if success
		// _, err = s.Write([]byte("ACK"))
		// if err != nil {
		// 	log.Printf("failed to write: %v\n", err)
		// }
	}
}
