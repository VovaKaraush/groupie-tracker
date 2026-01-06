package main

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

func fetchArtists() []Artist {
	var artists []Artist
	err := GetJson(baseURL+"/artists", &artists)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return artists
}

func fetchLocations() map[string][]string {
	var locations map[string][]string
	err := GetJson(baseURL+"/locations", &locations)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return locations
}

func fetchDates() map[string][]string {
	var dates map[string][]string
	err := GetJson(baseURL+"/dates", &dates)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return dates
}

func fetchRelations() map[string]Relation {
	var relations map[string]Relation
	err := GetJson(baseURL+"/relations", &relations)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return relations
}
