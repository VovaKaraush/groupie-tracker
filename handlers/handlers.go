package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"groupie-tracker/filters"
	"groupie-tracker/modules"
)

// Data holds all API data loaded at startup (artists, locations, dates, relations)
var Data *modules.GroupieData

// SortData initializes the global Data variable with API data
func SortData(g *modules.GroupieData) { Data = g }

// HomePageData contains all data needed to render the home.html template
// with search results and filter options
type HomePageData struct {
	// Filters stores raw string values to repopulate form fields
	// (keeps "0" from displaying in number inputs)
	Filters      FiltersView
	AllLocations []string
	Artists      []modules.Artist
}

// FiltersView contains the raw filter values from URL query parameters
// All values are stored as strings to preserve form input state
type FiltersView struct {
	Q string

	CreationDateMin string
	CreationDateMax string
	FirstAlbumMin   string
	FirstAlbumMax   string
	MembersMin      string
	MembersMax      string
	Locations       string // comma-separated city names from form
}

// homeTmpl is the main home page template with custom template functions
var homeTmpl = template.Must(template.New("home.html").Funcs(template.FuncMap{
	// intOrEmpty converts 0 to empty string, otherwise returns the integer
	"intOrEmpty": func(n int) string {
		if n == 0 {
			return ""
		}
		return strconv.Itoa(n)
	},
	// contains checks if a string is in a string slice
	"contains": func(list []string, item string) bool {
		for _, v := range list {
			if v == item {
				return true
			}
		}
		return false
	},
	// formatLocation converts "new_york-usa" to "New York - USA"
	"formatLocation": func(location string) string {
		// Replace underscores with spaces
		location = strings.ReplaceAll(location, "_", " ")
		// Capitalize each word
		words := strings.Fields(location)
		for i, word := range words {
			if len(word) > 0 {
				// Check if word contains a dash (country separator)
				if strings.Contains(word, "-") {
					// Split on dash and capitalize both parts
					parts := strings.Split(word, "-")
					for j, part := range parts {
						if len(part) > 0 {
							parts[j] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
						}
					}
					words[i] = strings.Join(parts, " - ")
				} else {
					words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
				}
			}
		}
		return strings.Join(words, " ")
	},
}).ParseFiles("web/home.html"))

// Concert represents a single concert event with location and dates
type Concert struct {
	Location string   // raw location string (e.g., "Paris-France")
	City     string   // extracted city name
	Country  string   // extracted country name
	Dates    []string // list of concert dates
}

// ArtistPageData contains all data needed to render the artist detail page
type ArtistPageData struct {
	Artist      modules.Artist
	ConcertList []Concert
}

// artistTmpl is the artist detail page template with custom formatting functions
var artistTmpl = template.Must(template.New("artist.html").Funcs(template.FuncMap{
	// formatLocation converts "new_york" to "New York"
	"formatLocation": func(location string) string {
		// Replace underscores with spaces
		location = strings.ReplaceAll(location, "_", " ")
		// Capitalize each word
		words := strings.Fields(location)
		for i, word := range words {
			if len(word) > 0 {
				words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
			}
		}
		return strings.Join(words, " ")
	},
	// googleMapsURL generates a Google Maps search link for a city and country
	"googleMapsURL": func(city, country string) string {
		// Create Google Maps search URL for city and country
		query := city
		if country != "" {
			query = city + " " + country
		}
		// URL encode the query by replacing spaces with +
		query = strings.ReplaceAll(query, " ", "+")
		return "https://www.google.com/maps/search/" + query
	},
}).ParseFiles("web/artist.html"))

// HomeHandler renders the main artist gallery page with search and filter functionality
// It processes query parameters, applies filters, and returns matching artists
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if Data == nil {
		http.Error(w, "Data not loaded yet", http.StatusInternalServerError)
		return
	}

	// Extract and parse all filter parameters from URL query string
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	creationDateMinRaw := strings.TrimSpace(r.URL.Query().Get("creationDateMin"))
	creationDateMaxRaw := strings.TrimSpace(r.URL.Query().Get("creationDateMax"))
	firstAlbumMinRaw := strings.TrimSpace(r.URL.Query().Get("firstAlbumMin"))
	firstAlbumMaxRaw := strings.TrimSpace(r.URL.Query().Get("firstAlbumMax"))
	membersMinRaw := strings.TrimSpace(r.URL.Query().Get("membersMin"))
	membersMaxRaw := strings.TrimSpace(r.URL.Query().Get("membersMax"))

	creationDateMin, _ := strconv.Atoi(creationDateMinRaw)
	creationDateMax, _ := strconv.Atoi(creationDateMaxRaw)
	firstAlbumMin, _ := strconv.Atoi(firstAlbumMinRaw)
	firstAlbumMax, _ := strconv.Atoi(firstAlbumMaxRaw)
	membersMin, _ := strconv.Atoi(membersMinRaw)
	membersMax, _ := strconv.Atoi(membersMaxRaw)

	// Parse location query - support both checkbox array and comma-separated formats
	locations := r.URL.Query()["locations"]
	// If only one location value with comma, treat it as comma-separated
	if len(locations) == 1 && strings.Contains(locations[0], ",") {
		locations = splitLocations(locations[0])
	}

	searchQuery := filters.SearchQuery{Query: q}
	criteria := filters.FilterCriteria{
		CreationDateMin: creationDateMin,
		CreationDateMax: creationDateMax,
		FirstAlbumMin:   firstAlbumMin,
		FirstAlbumMax:   firstAlbumMax,
		MembersMin:      membersMin,
		MembersMax:      membersMax,
		Locations:       locations,
	}

	// If no filters applied, show all artists; otherwise apply search and filters
	noFilters := q == "" &&
		creationDateMin == 0 && creationDateMax == 0 &&
		firstAlbumMin == 0 && firstAlbumMax == 0 &&
		membersMin == 0 && membersMax == 0 &&
		len(locations) == 0

	var results []modules.Artist
	if noFilters {
		results = Data.Artists
	} else {
		results = filters.SearchWithFilters(Data, searchQuery, criteria)
	}

	page := HomePageData{
		Filters: FiltersView{
			Q:               q,
			CreationDateMin: creationDateMinRaw,
			CreationDateMax: creationDateMaxRaw,
			FirstAlbumMin:   firstAlbumMinRaw,
			FirstAlbumMax:   firstAlbumMaxRaw,
			MembersMin:      membersMinRaw,
			MembersMax:      membersMaxRaw,
			Locations:       strings.Join(locations, ","),
		},
		AllLocations: uniqueLocations(Data),
		Artists:      results,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := homeTmpl.Execute(w, page); err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// SearchHandler returns filtered artists as JSON for AJAX requests
// Supports both comma-separated and multiple location parameters
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query().Get("q")
	creationDateMin, _ := strconv.Atoi(r.URL.Query().Get("creationDateMin"))
	creationDateMax, _ := strconv.Atoi(r.URL.Query().Get("creationDateMax"))
	firstAlbumMin, _ := strconv.Atoi(r.URL.Query().Get("firstAlbumMin"))
	firstAlbumMax, _ := strconv.Atoi(r.URL.Query().Get("firstAlbumMax"))
	membersMin, _ := strconv.Atoi(r.URL.Query().Get("membersMin"))
	membersMax, _ := strconv.Atoi(r.URL.Query().Get("membersMax"))
	// Support both text field (comma-separated) and multiple query params formats
	// - text field: "locations=Paris, Lyon"
	// - multiple params: "locations=Paris&locations=Lyon"
	locations := r.URL.Query()["locations"]
	if len(locations) == 1 {
		locations = splitLocations(locations[0])
	}

	searchQuery := filters.SearchQuery{Query: query}
	criteria := filters.FilterCriteria{
		CreationDateMin: creationDateMin,
		CreationDateMax: creationDateMax,
		FirstAlbumMin:   firstAlbumMin,
		FirstAlbumMax:   firstAlbumMax,
		MembersMin:      membersMin,
		MembersMax:      membersMax,
		Locations:       locations,
	}

	results := filters.SearchWithFilters(Data, searchQuery, criteria)
	_ = json.NewEncoder(w).Encode(results)
}

// ArtistHandler is a fallback handler (legacy) that displays all artists
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	if Data == nil {
		http.Error(w, "Data not loaded yet", http.StatusInternalServerError)
		return
	}
	// Display all artists without filters
	page := HomePageData{
		Filters:      FiltersView{},
		AllLocations: uniqueLocations(Data),
		Artists:      Data.Artists,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = homeTmpl.Execute(w, page)
}

// splitLocations parses a comma-separated string into a trimmed string slice
// Returns nil if input is empty
func splitLocations(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// uniqueLocations extracts and sorts all unique city names from the API locations data
func uniqueLocations(data *modules.GroupieData) []string {
	set := map[string]struct{}{}
	for _, l := range data.Locations {
		for _, city := range l.Locations {
			set[city] = struct{}{}
		}
	}
	out := make([]string, 0, len(set))
	for city := range set {
		out = append(out, city)
	}
	sort.Strings(out)
	return out
}

// DetailHandler displays a detailed page for a single artist including members and concerts
// Artist ID is passed as a URL query parameter (?id=1)
func DetailHandler(w http.ResponseWriter, r *http.Request) {
	if Data == nil {
		http.Error(w, "Data not loaded yet", http.StatusInternalServerError)
		return
	}

	// Extract and validate artist ID from query parameters
	idRaw := strings.TrimSpace(r.URL.Query().Get("id"))
	if idRaw == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(idRaw)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Find the artist with matching ID in global data
	var artist *modules.Artist
	for i := range Data.Artists {
		if Data.Artists[i].ID == id {
			artist = &Data.Artists[i]
			break
		}
	}

	if artist == nil {
		http.NotFound(w, r)
		return
	}

	// Parse and extract concert information from Relations data
	concerts := extractConcerts(id)

	page := ArtistPageData{
		Artist:      *artist,
		ConcertList: concerts,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := artistTmpl.Execute(w, page); err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// extractConcerts builds Concert entries for an artist by parsing Relations data
// Locations are assumed to be in format "City-Country"
// The API groups concert dates by location
func extractConcerts(artistID int) []Concert {
	var concerts []Concert

	// Look up relation data for this artist (maps locations to concert dates)
	var datesLocations map[string]interface{}
	for _, rel := range Data.Relations {
		if rel.ID == artistID {
			datesLocations = rel.DatesLocations
			break
		}
	}

	if datesLocations == nil {
		return concerts
	}

	// Iterate through each location and its associated concert dates
	for location, dateInterface := range datesLocations {
		dates := []string{}

		// Type assert dates from interface{} to []interface{}
		if dateList, ok := dateInterface.([]interface{}); ok {
			for _, d := range dateList {
				if dateStr, ok := d.(string); ok {
					dates = append(dates, dateStr)
				}
			}
		}

		if len(dates) > 0 {
			// Extract city and country from "City-Country" format
			city := location
			country := ""
			if idx := strings.Index(location, "-"); idx != -1 {
				city = strings.TrimSpace(location[:idx])
				country = strings.TrimSpace(location[idx+1:])
			}

			concerts = append(concerts, Concert{
				Location: location,
				City:     city,
				Country:  country,
				Dates:    dates,
			})
		}
	}

	// Sort concerts alphabetically by location for consistent display
	sort.Slice(concerts, func(i, j int) bool {
		return concerts[i].Location < concerts[j].Location
	})

	return concerts
}
