package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/SmartHouse/1.0.2/",
		Index,
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
		"GetLightState",
		"GET",
		"/SmartHouse/1.0.2/lights/{lightID}",
		GetLightState,
	},

	Route{
		"GetLights",
		"GET",
		"/SmartHouse/1.0.2/lights",
		GetLights,
	},

	Route{
		"SetLightState",
		"PUT",
		"/SmartHouse/1.0.2/lights/{lightID}/{state}",
		SetLightState,
	},

	Route{
		"SetLuminosityThreshhold",
		"PUT",
		"/SmartHouse/1.0.2/lights/{lightID}/settings",
		SetLuminosityThreshhold,
	},

	Route{
		"SettingsLight",
		"GET",
		"/SmartHouse/1.0.2/lights/{lightID}/settings",
		SettingsLight,
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
		"SetHomeSettings",
		"PUT",
		"/SmartHouse/1.0.2/settings/home/",
		SetHomeSettings,
	},
}

type StatusResponse struct {
	Message string `json:"message,omitempty"`
}
