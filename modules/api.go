package modules

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var baseURL = "https://groupietrackers.herokuapp.com/api"

var httpClient = &http.Client{Timeout: 10 * time.Second}

func GetJson(url string, data interface{}) error {
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		resp, err := httpClient.Get(url)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(attempt) * 200 * time.Millisecond)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 8*1024))
			resp.Body.Close()
			lastErr = fmt.Errorf("unexpected status %d for %s: %s", resp.StatusCode, url, string(body))
			time.Sleep(time.Duration(attempt) * 200 * time.Millisecond)
			continue
		}

		decErr := json.NewDecoder(resp.Body).Decode(data)
		resp.Body.Close()
		if decErr == nil {
			return nil
		}

		lastErr = decErr
		time.Sleep(time.Duration(attempt) * 200 * time.Millisecond)
	}
	return lastErr
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
