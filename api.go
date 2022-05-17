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
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`     // nil | string
	Previous interface{} `json:"previous"` // nil | string
	Results  []T         `json:"results"`
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
	Id        int         `json:"id"`
	Name      string      `json:"name"`
	Height    int         `json:"height"`
	Created   time.Time   `json:"created"`
	Edited    time.Time   `json:"edited"`
	Homeworld int         `json:"homeworld"` // the ID only
	Mass      interface{} `json:"mass"`      // int or unknown
}

// Used to tackle overfetching
type planetDTO struct {
	Name       string `json:"name"`
	Diameter   int    `json:"diameter"`
	Climate    string `json:"climate"`
	Population int    `json:"population"`
}

// Would return 2 for a URL such as "https://swapi.dev/api/planets/2/"
func getResourceIDFromURL(url string) int {
	urlSplit := strings.Split(url, "/")
	id := urlSplit[len(urlSplit)-2]
	intID, _ := strconv.Atoi(id)
	return intID
}

func numericStringOrUnknownToIntOrNil(s string) interface{} {
	if s == "unknown" {
		return nil
	} else {
		number, _ := strconv.Atoi(s)
		return number
	}
}

func swapiPersonToPerson(swapiPerson swapiPersonDTO) personDTO {
	height, _ := strconv.Atoi(swapiPerson.Height)

	return personDTO{
		Id:        getResourceIDFromURL(swapiPerson.URL),
		Name:      swapiPerson.Name,
		Height:    height,
		Mass:      numericStringOrUnknownToIntOrNil(swapiPerson.Mass),
		Created:   swapiPerson.Created,
		Edited:    swapiPerson.Edited,
		Homeworld: getResourceIDFromURL(swapiPerson.Homeworld),
	}
}

func swapiPlanetToPlanet(swapiPlanet swapiPlanetDTO) planetDTO {
	diameter, _ := strconv.Atoi(swapiPlanet.Diameter)
	population, _ := strconv.Atoi(swapiPlanet.Population)

	return planetDTO{
		Name:       swapiPlanet.Name,
		Diameter:   diameter,
		Climate:    swapiPlanet.Climate,
		Population: population,
	}
}

func getPeopleFromPage(page int) swapiPeopleReponse {
	resp, err := http.Get(strings.Join([]string{"http://swapi.dev/api/people?page=", strconv.Itoa(page)}, ""))

	if err != nil {
		// TODO infinite retry could totally backfire
		return getPeopleFromPage(page)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)

		var unmarshalled swapiPeopleReponse
		unmarshallingError := json.Unmarshal(body, &unmarshalled)

		if unmarshallingError != nil {
			fmt.Println("unmarshalling error thrown")
			fmt.Println(unmarshallingError)
		}

		return unmarshalled
	}
}

// Fetch and store all people beforehand
func getAllPeopleFromSwapi() []swapiPersonDTO {
	firstRes := getPeopleFromPage(1)

	results := make([]swapiPersonDTO, firstRes.Count)
	for i, v := range firstRes.Results {
		results[i] = v
	}

	pages := int(math.Ceil(float64(firstRes.Count) / 10))

	for page := 2; page <= pages; page++ {
		for i, v := range getPeopleFromPage(page).Results {
			results[i+page*10-10] = v
		}
	}

	return results
}

// TODO issue -> it is sync
// TODO receive pageSize arg
func handleGetPeople(writer http.ResponseWriter, req *http.Request) {
	resp, err := http.Get("http://swapi.dev/api/people")
	writer.Header().Add("Content-Type", "application/json")

	if err == nil {
		body, err := ioutil.ReadAll(resp.Body)

		var unmarshalled swapiPeopleReponse
		unmarshallingError := json.Unmarshal(body, &unmarshalled)

		if unmarshallingError != nil {
			fmt.Println("unmarshalling error thrown")
			fmt.Println(unmarshallingError)
		}

		results := make([]personDTO, len(unmarshalled.Results))
		for i, v := range unmarshalled.Results {
			results[i] = swapiPersonToPerson(v)
		}
		resultsJSON, _ := json.Marshal(results)

		if err != nil {
			writer.WriteHeader(400)
		} else {
			writer.Write(resultsJSON)
		}

	} else {
		writer.WriteHeader(400)
	}
}

func main() {
	ppl := getAllPeopleFromSwapi()
	fmt.Print(ppl)
	// TODO fetch all people beforeHand
	// http.HandleFunc("/people", handleGetPeople)
	// http.ListenAndServe(":8080", nil)
}
