package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type App struct {
	Env   string
	Games []Game
}

type Game struct {
	ID   string
	Name string
	Date string
}

// func GetGames() string {
func GetGames() []Game {
	backend := "http://localhost"
	if b := os.Getenv("BACKEND_URL"); b != "" {
		backend = b
	}

	url := backend + "/games"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var games []Game
	err = json.Unmarshal([]byte(body), &games)
	if err != nil {
		log.Fatal(err)
	}

	return games
}

func IndexHandler(response http.ResponseWriter, request *http.Request) {
	tmplt := template.New("index.html")
	tmplt, _ = tmplt.ParseFiles("templates/index.html")

	env := os.Getenv("ENV_NAME")
	games := GetGames()
	r := App{Env: env,
		Games: games}

	tmplt.Execute(response, r)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/status", HealthCheck)
	http.ListenAndServe(":80", nil)
}
