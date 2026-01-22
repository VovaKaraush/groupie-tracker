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

func FetchLocations() []Location {
	var locResp LocationsResponse
	err := GetJson(baseURL+"/locations", &locResp)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return locResp.Index
}

func FetchDates() []DateInfo {
	var datesResp DatesResponse
	err := GetJson(baseURL+"/dates", &datesResp)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return datesResp.Index
}

func FetchRelations() []Relation {
	var relResp RelationResponse
	err := GetJson(baseURL+"/relation", &relResp)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return relResp.Index
}
