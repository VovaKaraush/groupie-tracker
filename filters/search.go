package filters

import (
	"groupie-tracker/modules"
	"strings"
)

type SearchQuery struct {
	Query string
}

// ça combine la recherche avec les critères de filtrage
type SearchFilters struct {
	SearchQuery    SearchQuery
	FilterCriteria FilterCriteria
}

func SearchArtists(data *modules.GroupieData, query string) []modules.Artist {
	var results []modules.Artist
	query = strings.ToLower(strings.TrimSpace(query))

	if query == "" {
		return data.Artists
	}

	for _, artist := range data.Artists {
		// Rechercher dans le nom de l'artiste
		if strings.Contains(strings.ToLower(artist.Name), query) {
			results = append(results, artist)
			continue
		}

		// Rechercher dans les noms des membres
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				results = append(results, artist)
				break
			}
		}
	}

	return results
}

// ça effectue un recherche combinée avec les filtre
func SearchWithFilters(data *modules.GroupieData, search SearchQuery, criteria FilterCriteria) []modules.Artist {
	if data == nil || len(data.Artists) == 0 {
		return []modules.Artist{}
	}

	searchResults := SearchArtists(data, search.Query)

	// Si pas de critère de filtrage, retourner les résultats de recherche
	if criteria.CreationDateMin == 0 && criteria.CreationDateMax == 0 &&
		criteria.FirstAlbumMin == 0 && criteria.FirstAlbumMax == 0 &&
		criteria.MembersMin == 0 && criteria.MembersMax == 0 &&
		len(criteria.Locations) == 0 {
		return searchResults
	}

	tempData := &modules.GroupieData{
		Artists:   searchResults,
		Locations: data.Locations,
		Dates:     data.Dates,
		Relations: data.Relations,
	}

	// on appliquer les filtres sur les résultat de recherche
	return ApplyFilters(tempData, criteria)
}

func ValidateSearchQuery(query string) bool {
	return len(strings.TrimSpace(query)) > 0
}
