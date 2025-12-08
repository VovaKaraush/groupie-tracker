package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var baseURL = "https://groupietrackers.herokuapp.com/api"
var client = &http.Client{}

type Artists struct {
	Artists []Artist
}

type Artist struct { // mostly would be groups
	id           int      `json:id`
	image        string   `json:image`
	name         string   `json:name`
	members      []string `json:members`
	creationDate int      `json:creationDate`
	firstAlbum   string   `json:firstAlbum`
	locations    string   `json:locations`
	concertDates string   `json:concertDates`
	relations    string   `json:relations`
}

type Locations struct {
	Locations []Location
}

type Location struct {
	id        int    `json:id`
	locations string `json:locations`
	dates     string `json:dates`
}

type Dates struct {
	Dates []Date
}

type Date struct {
	id    int    `json:id`
	dates string `json:dates`
}

type Relations struct {
	Relations []Relation
}

type Relation struct {
	id             int                 `json:id`
	datesLocations map[string][]string `json:datesLocations`
}

func main() {
	artists := fetchArtists()
	locations := fetchLocations()
	dates := fetchDates()
	relations := fetchRelations()
	log.Println("Fetches : ")
	log.Printf("Fetched %d artists", len(artists.Artists))
	log.Printf("Fetched %d locations", len(locations.Locations))
	log.Printf("Fetched %d dates", len(dates.Dates))
	log.Printf("Fetched %d relations", len(relations.Relations))
}

func GetJson(url string, data interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(data)
}

func fetchArtists() Artists {
	var artists Artists
	err := GetJson("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		fmt.Printf("%v", err.Error())
	}
	return artists
}

func fetchLocations() Locations {
	var locations Locations
	err := GetJson("https://groupietrackers.herokuapp.com/api/locations", &locations)
	if err != nil {
		fmt.Printf("%v", err.Error())
	}
	return locations
}

func fetchDates() Dates {
	var dates Dates
	err := GetJson("https://groupietrackers.herokuapp.com/api/dates", &dates)
	if err != nil {
		fmt.Printf("%v", err.Error())
	}
	return dates
}

func fetchRelations() Relations {
	var relations Relations
	err := GetJson("https://groupietrackers.herokuapp.com/api/relations", &relations)
	if err != nil {
		fmt.Printf("%v", err.Error())
	}
	return relations
}
