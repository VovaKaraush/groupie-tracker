package main

import (
	"groupie-tracker/modules"
	"log"
)

func main() {
	// Création d'une structure pour regrouper toutes les données
	data := modules.GroupieData{
		Artists:   modules.FetchArtists(),
		Locations: modules.FetchLocations(),
		Dates:     modules.FetchDates(),
		Relations: modules.FetchRelations(),
	}

	log.Println("Fetches : ")
	log.Printf("Fetched %d artists", len(data.Artists))
	log.Printf("Fetched %d locations", len(data.Locations))
	log.Printf("Fetched %d dates", len(data.Dates))
	log.Printf("Fetched %d relations", len(data.Relations))

}
