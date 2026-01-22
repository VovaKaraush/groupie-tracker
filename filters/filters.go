// Les filtres devraient permettre de filtrer par date de création, date de l'album, nommbre de membres et par localisation.
// les filtres devraient aussi être de type "range filter" et "checkbox filter".
package filters

import (
	"fmt"
	"groupie-tracker/modules"
)

// FilterCriteria définit les critères de filtrage
type FilterCriteria struct {
	CreationDateMin int
	CreationDateMax int
	FirstAlbumMin   int
	FirstAlbumMax   int
	MembersMin      int
	MembersMax      int
	Locations       []string
}

// ApplyFilters applique les filtres sur les données et retourne les artistes filtrés
func ApplyFilters(data *modules.GroupieData, criteria FilterCriteria) []modules.Artist {
	var filtered []modules.Artist

	for _, artist := range data.Artists {
		// Filtrer par date de création
		if criteria.CreationDateMin > 0 && artist.CreationDate < criteria.CreationDateMin {
			continue
		}
		if criteria.CreationDateMax > 0 && artist.CreationDate > criteria.CreationDateMax {
			continue
		}

		// Filtrer par nombre de membres
		if criteria.MembersMin > 0 && len(artist.Members) < criteria.MembersMin {
			continue
		}
		if criteria.MembersMax > 0 && len(artist.Members) > criteria.MembersMax {
			continue
		}

		// filtrer par localisation (via data.Locations par ID d'artiste)
		if len(criteria.Locations) > 0 {
			matched := false
			// Récupérer les localisations de l'artiste à partir de data.Locations
			var artistLocations []string
			for _, loc := range data.Locations {
				if loc.ID == artist.ID {
					artistLocations = loc.Locations
					break
				}
			}
			for _, loc := range criteria.Locations {
				for _, artistLoc := range artistLocations {
					if loc == artistLoc {
						matched = true
						break
					}
				}
				if matched {
					break
				}
			}
			if !matched {
				continue
			}
		}

		// Filtrer par date du premier album ( format DD-MM-YYYY)
		if criteria.FirstAlbumMin > 0 {
			albumYear := 0
			_, err := fmt.Sscanf(artist.FirstAlbum, "%d-", &albumYear)
			if err == nil && albumYear < criteria.FirstAlbumMin {
				continue
			}
		}
		if criteria.FirstAlbumMax > 0 {
			albumYear := 0
			_, err := fmt.Sscanf(artist.FirstAlbum, "%d-", &albumYear)
			if err == nil && albumYear > criteria.FirstAlbumMax {
				continue
			}
		}

		// Si l'artiste passe tous les filtres, l'ajouter à la liste filtrée
		filtered = append(filtered, artist)
	}

	return filtered
}
