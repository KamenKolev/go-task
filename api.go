package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type swapiPersonDTO struct {
	Name      string    `json:"name"`
	Height    string    `json:"height"`
	Mass      string    `json:"mass"`
	HairColor string    `json:"hair_color"`
	SkinColor string    `json:"skin_color"`
	EyeColor  string    `json:"eye_color"`
	BirthYear string    `json:"birth_year"`
	Gender    string    `json:"gender"`
	Homeworld string    `json:"homeworld"`
	Films     []string  `json:"films"`
	Species   []string  `json:"species"`
	Vehicles  []string  `json:"vehicles"`
	Starships []string  `json:"starships"`
	Created   time.Time `json:"created"`
	Edited    time.Time `json:"edited"`
	URL       string    `json:"url"`
}

type swapiPersonsReponse struct {
	Count    int16            `json:"count"`
	Next     interface{}      `json:"next"`     // nil | string
	Previous interface{}      `json:"previous"` // nil | string
	Results  []swapiPersonDTO `json:"results"`
}

type swapiPlanetDTO struct {
	Name           string    `json:"name"`
	RotationPeriod string    `json:"rotation_period"`
	OrbitalPeriod  string    `json:"orbital_period"`
	Diameter       string    `json:"diameter"`
	Climate        string    `json:"climate"`
	Gravity        string    `json:"gravity"`
	Terrain        string    `json:"terrain"`
	SurfaceWater   string    `json:"surface_water"`
	Population     string    `json:"population"`
	Residents      []string  `json:"residents"`
	Films          []string  `json:"films"`
	Created        time.Time `json:"created"`
	Edited         time.Time `json:"edited"`
	URL            string    `json:"url"`
}

type swapiPlanetsReponse struct {
	Count    int16            `json:"count"`
	Next     interface{}      `json:"next"`     // nil | string
	Previous interface{}      `json:"previous"` // nil | string
	Results  []swapiPlanetDTO `json:"results"`
}

// Used to tackle overfetching
type personDTO struct {
	Id        int16       `json:"id"`
	Name      string      `json:"name"`
	Height    int         `json:"height"`
	Created   string      `json:"created"`
	Edited    string      `json:"edited"`
	Homeworld int16       `json:"homeworld"`
	Mass      interface{} `json:"mass"` // int or unknown
}

// Used to tackle overfetching
type planetDTO struct {
	Name       string `json:"name"`
	Diameter   int    `json:"diameter"`
	Climate    string `json:"climate"`
	Population int    `json:"population"`
}

func hello(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, "Yo!")
}

func getPeople(writer http.ResponseWriter, req *http.Request) {
	resp, err := http.Get("http://swapi.dev/api/people")
	writer.Header().Add("Content-Type", "application/json")

	if err == nil {
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			writer.WriteHeader(400)
		} else {
			writer.Write(body)

		}

	} else {
		writer.WriteHeader(400)
	}
}

func main() {
	http.HandleFunc("/people", getPeople)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Listening at port 8080")
}
