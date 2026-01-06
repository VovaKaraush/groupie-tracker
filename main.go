package main

import (
	"groupie-tracker/modules"
	"log"
)

func main() {
	artists := modules.FetchArtists()
	locations := modules.FetchLocations()
	dates := modules.FetchDates()
	relations := modules.FetchRelations()

	log.Println("Fetches : ")
	log.Printf("Fetched %d artists", len(artists))
	log.Printf("Fetched %d locations", len(locations))
	log.Printf("Fetched %d dates", len(dates))
	log.Printf("Fetched %d relations", len(relations))

}
