package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type personDTO struct {
	name       string
	height     string
	mass       string
	hair_color string
	skin_color string // can be multiple, comma separated
	eye_color  string
	birth_year string
	gender     string // enum
	homeworld  string
	films      []string
	species    []string
	vehicles   []string
	starships  []string
	created    string
	edited     string
	url        string
}

type planetDTO struct {
	name            string
	rotation_period string
	orbital_period  string
	diameter        string
	climate         string
	gravity         string
	terrain         string
	surface_water   string
	population      string
	residents       []string
	films           []string
	created         string
	edited          string
	url             string
}

type person struct {
	id        int16
	name      string
	height    int
	created   string
	edited    string
	homeworld int16 // id
	mass      int   // or unknown
}

type planet struct {
	name       string
	diameter   string
	climate    string
	population int
}

func hello(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, "Yo!")
}

func getPeople(writer http.ResponseWriter, req *http.Request) {
	resp, err := http.Get("http://swapi.dev/api/people")

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
