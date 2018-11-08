package main

import (
    "net/http"
    "net/http/httptest"
	"testing"
    "strings"
    "github.com/gorilla/mux"
)

func TestAllGames(t *testing.T) {
    games = append(games, Games{ID: "0", Name: "test", Date: "data"})

	req, err := http.NewRequest("GET", "/games", nil)

    if err != nil {
		t.Errorf("An error occurred. %v", err)
	}

	rr := httptest.NewRecorder()

	r := mux.NewRouter()

	r.HandleFunc("/games", AllGames).Methods("GET")

	r.ServeHTTP(rr, req)

    expected := `{"ID":"0","Name":"test","Date":"data"}`
  
    assert := strings.Contains(rr.Body.String(), expected)
    if !assert {
        t.Errorf("Unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}

func TestGetDate(t *testing.T) {
    str := "10/03/2019"
    result := GetDate(str)
    if result != str {
        t.Errorf("wrong string returned %v", result)
    }
}

func TestArticleHandlerWithAnInvalidPost(t *testing.T) {
    games = append(games, Games{ID: "0", Name: "test", Date: "data"})

	req, err := http.NewRequest("GET", "/game/0", nil)

    if err != nil {
		t.Errorf("An error occurred. %v", err)
	}

	rr := httptest.NewRecorder()

	r := mux.NewRouter()

	r.HandleFunc("/game/{id}", GetGame).Methods("GET")

	r.ServeHTTP(rr, req)

    expected := `{"ID":"0","Name":"test","Date":"data"}`
    
    assert := strings.Contains(rr.Body.String(), expected)
    if !assert {
        t.Errorf("Unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/hc", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	HealthCheck(res, req)
	if res.Code != 200 {
		t.Fatal("Response code is not 200")
	}
}