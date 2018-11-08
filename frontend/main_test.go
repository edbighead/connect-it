package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

type Games []Game

var g = Games{
	Game{
		"Some",
		"Test",
		"Data",
	},
	Game{
		"Some2",
		"Test2",
		"Data2",
	},
}

func TestIndexHandler(t *testing.T) {

	defer gock.Off()
	gamesJson, err := json.Marshal(g)
	gock.New("http://localhost").
		Get("/games").
		Reply(200).
		JSON(gamesJson)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	IndexHandler(res, req)

	exp := "COMING SOON"
	act := res.Body.String()

	assert := strings.Contains(act, exp)

	if !assert {
		t.Fatalf("Expected string %s in %s", exp, act)
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
