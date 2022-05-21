package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type swapiMultipleResourcesResponse[T any] struct {
	Count    int `json:"count"`
	Next     any `json:"next"`     // nil | string
	Previous any `json:"previous"` // nil | string
	Results  []T `json:"results"`
}

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

type swapiPeopleReponse = swapiMultipleResourcesResponse[swapiPersonDTO]

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

type swapiPlanetsReponse = swapiMultipleResourcesResponse[swapiPlanetDTO]

// Used to tackle overfetching
type personDTO struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Height    any       `json:"height"` // float or nil
	Created   time.Time `json:"created"`
	Edited    time.Time `json:"edited"`
	Homeworld any       `json:"homeworld"` // the ID only. Could be nil if the planet is unknown. In this case, this is planet 28
	Mass      any       `json:"mass"`      // float or nil
}

// Used to tackle overfetching
type planetDTO struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Diameter   any    `json:"diameter"` // float or nil
	Climate    string `json:"climate"`
	Population any    `json:"population"` // float or nil
}

// Would return 2 for a URL such as "https://swapi.dev/api/planets/2/"
func getResourceIDFromURL(url string) (int, error) {
	urlSplit := strings.Split(url, "/")
	id := urlSplit[len(urlSplit)-2]
	intID, stringConversionError := strconv.Atoi(id)
	return intID, stringConversionError
}

func numericStringOrUnknownToFloatOrNil(s string) (any, error) {
	if s == "unknown" {
		return nil, nil
	}

	// The API uses commas to signify thousands. They don't play well with ParseFloat
	s = strings.ReplaceAll(s, ",", "")

	number, convError := strconv.ParseFloat(s, 64)
	if convError != nil {
		return nil, convError
	}

	return number, nil
}

func swapiPersonToPerson(swapiPerson swapiPersonDTO) (personDTO, error) {
	height, heightConversionError := numericStringOrUnknownToFloatOrNil(swapiPerson.Height)
	if heightConversionError != nil {
		return personDTO{}, heightConversionError
	}
	id, idConversionError := getResourceIDFromURL(swapiPerson.URL)
	if idConversionError != nil {
		return personDTO{}, idConversionError
	}

	var homeworld any
	homeworld, homeworldConversionError := getResourceIDFromURL(swapiPerson.Homeworld)
	if homeworldConversionError != nil {
		return personDTO{}, homeworldConversionError
	}
	mass, massConvError := numericStringOrUnknownToFloatOrNil(swapiPerson.Mass)
	if massConvError != nil {
		return personDTO{}, massConvError
	}

	// Planet 28 is "unknown" and has zero other useful info. Better return null instead
	if homeworld == 28 {
		homeworld = nil
	}

	return personDTO{
		Id:        id,
		Name:      swapiPerson.Name,
		Height:    height,
		Mass:      mass,
		Created:   swapiPerson.Created,
		Edited:    swapiPerson.Edited,
		Homeworld: homeworld,
	}, nil
}

func swapiPlanetToPlanet(swapiPlanet swapiPlanetDTO) (planetDTO, error) {
	diameter, diameterConvError := numericStringOrUnknownToFloatOrNil(swapiPlanet.Diameter)
	if diameterConvError != nil {
		return planetDTO{}, diameterConvError
	}

	population, populationConvError := numericStringOrUnknownToFloatOrNil(swapiPlanet.Population)
	if populationConvError != nil {
		return planetDTO{}, populationConvError
	}
	id, idError := getResourceIDFromURL(swapiPlanet.URL)
	if idError != nil {
		return planetDTO{}, idError
	}

	return planetDTO{
		Id:         id,
		Name:       swapiPlanet.Name,
		Diameter:   diameter,
		Climate:    swapiPlanet.Climate,
		Population: population,
	}, nil
}

func getPeopleFromPage(page int) (swapiPeopleReponse, error) {
	resp, err := http.Get(strings.Join([]string{"http://swapi.dev/api/people?page=", strconv.Itoa(page)}, ""))

	if err != nil {
		// TODO infinite retry could totally backfire
		fmt.Print("retrying getPeopleFromPage")
		return getPeopleFromPage(page)
	} else {
		body, readingError := ioutil.ReadAll(resp.Body)
		var unmarshalled swapiPeopleReponse

		if readingError != nil {
			fmt.Println("readingError error thrown")
			fmt.Println(readingError)
			return unmarshalled, readingError
		}

		unmarshallingError := json.Unmarshal(body, &unmarshalled)

		if unmarshallingError != nil {
			fmt.Println("unmarshalling error thrown")
			fmt.Println(unmarshallingError)
			return unmarshalled, unmarshallingError
		}

		return unmarshalled, nil
	}
}

// Fetch and store all people beforehand
func getAllPeopleFromSwapi() []swapiPersonDTO {
	firstRes, _ := getPeopleFromPage(1)

	results := make([]swapiPersonDTO, firstRes.Count)
	for i, v := range firstRes.Results {
		results[i] = v
	}

	pages := int(math.Ceil(float64(firstRes.Count) / 10))

	for page := 2; page <= pages; page++ {
		res, error := getPeopleFromPage(page)
		if error != nil {
			fmt.Print(error)
		}
		for i, v := range res.Results {
			results[i+page*10-10] = v
		}
	}

	return results
}

// TODO dedupe
func getPlanetsFromPage(page int) swapiPlanetsReponse {
	resp, err := http.Get(strings.Join([]string{"http://swapi.dev/api/planets?page=", strconv.Itoa(page)}, ""))

	if err != nil {
		// TODO infinite retry could totally backfire
		return getPlanetsFromPage(page)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)

		var unmarshalled swapiPlanetsReponse
		unmarshallingError := json.Unmarshal(body, &unmarshalled)

		if unmarshallingError != nil {
			fmt.Println("unmarshalling error thrown")
			fmt.Println(unmarshallingError)
		}

		return unmarshalled
	}
}

func getAllPlanetsFromSwapi() []swapiPlanetDTO {
	firstRes := getPlanetsFromPage(1)

	results := make([]swapiPlanetDTO, firstRes.Count)
	for i, v := range firstRes.Results {
		results[i] = v
	}

	pages := int(math.Ceil(float64(firstRes.Count) / 10))

	for page := 2; page <= pages; page++ {
		for i, v := range getPlanetsFromPage(page).Results {
			results[i+page*10-10] = v
		}
	}

	return results
}

func handleGetPeople(storedPeople []swapiPersonDTO) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		writer.Header().Add("Access-Control-Allow-Origin", "*")

		results := make([]personDTO, len(storedPeople))
		for i, v := range storedPeople {
			person, err := swapiPersonToPerson(v)
			if err != nil {
				fmt.Println("swapiPersonToPerson error thrown", err, v)
				// TODO handle error
			}
			results[i] = person
		}
		resultsJSON, marshallingError := json.Marshal(results)

		if marshallingError != nil {
			fmt.Print(marshallingError)
		} else {
			writer.Write(resultsJSON)
		}
	}
}

func handleGetPlanets(storedPlanets []swapiPlanetDTO) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		writer.Header().Add("Access-Control-Allow-Origin", "*")

		results := make([]planetDTO, len(storedPlanets))
		for i, v := range storedPlanets {
			planet, err := swapiPlanetToPlanet(v)
			if err != nil {
				fmt.Println("swapiPlanetToPlanet error thrown", err, v)
				// TODO handle error
			}
			results[i] = planet
		}
		resultsJSON, marshallingError := json.Marshal(results)

		if marshallingError != nil {
			fmt.Print(marshallingError)
		} else {
			writer.Write(resultsJSON)
		}
	}
}

func main() {
	ppl := getAllPeopleFromSwapi()
	planets := getAllPlanetsFromSwapi()

	http.HandleFunc("/people", handleGetPeople(ppl))
	http.HandleFunc("/planets", handleGetPlanets(planets))
	http.ListenAndServe(":8080", nil)
}
