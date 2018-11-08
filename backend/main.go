package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Games struct {
	ID   string
	Name string
	Date string
}

var games = []Games{{ID: "1", Name: "Days Gone", Date: GetDate("02/22/2019")},
	{ID: "2", Name: "Last of Us: Part II", Date: GetDate("unknown")},
	{ID: "3", Name: "Death Stranding", Date: GetDate("unknown")}}

func AllGames(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	encoder.Encode(games)
}

func GetGame(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	encoder := json.NewEncoder(w)
	for _, item := range games {
		if item.ID == params["id"] {
			encoder.Encode(item)
			return
		}
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func GetDate(date string) string {
	if date != "unknown" {
		return date
	} else {
		return "TBA"
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/games", AllGames).Methods("GET")
	router.HandleFunc("/game/{id:[0-9]+}", GetGame).Methods("GET")
	router.HandleFunc("/status", HealthCheck).Methods("GET")
	log.Fatal(http.ListenAndServe(":80", router))
}
