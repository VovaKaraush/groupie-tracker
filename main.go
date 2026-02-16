package main

import (
	"fmt"
	"groupie-tracker/handlers"
	"groupie-tracker/modules"
	"log"
	"net/http"
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

	handlers.SortData(&data)

	// Serve static assets (CSS, images, etc.) from /web/.
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/search", handlers.SearchHandler)
	http.HandleFunc("/artist", handlers.ArtistHandler)

	fmt.Println("Serveur démarré sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
