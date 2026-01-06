package main

import (
	"log"
)

func main() {
	artists := fetchArtists()
	locations := fetchLocations()
	dates := fetchDates()
	relations := fetchRelations()

	log.Println("Fetches : ")
	log.Printf("Fetched %d artists", len(artists))
	log.Printf("Fetched %d locations", len(locations))
	log.Printf("Fetched %d dates", len(dates))
	log.Printf("Fetched %d relations", len(relations))

}
