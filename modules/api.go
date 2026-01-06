package modules

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var baseURL = "https://groupietrackers.herokuapp.com/api"

func GetJson(url string, data interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(data)
}

func FetchArtists() []Artist {
	var artists []Artist
	err := GetJson(baseURL+"/artists", &artists)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return artists
}

func FetchLocations() map[string][]string {
	var locations map[string][]string
	err := GetJson(baseURL+"/locations", &locations)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return locations
}

func FetchDates() map[string][]string {
	var dates map[string][]string
	err := GetJson(baseURL+"/dates", &dates)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return dates
}

func FetchRelations() map[string]Relation {
	var relations map[string]Relation
	err := GetJson(baseURL+"/relations", &relations)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return relations
}
