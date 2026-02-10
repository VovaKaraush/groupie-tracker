package handlers

import (
	"groupie-tracker/filters"
	"groupie-tracker/modules"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

var Data *modules.GroupieData

func SortData(g *modules.GroupieData) {
	Data = g
}

// PageData contient les artistes à afficher et les valeurs du formulaire
type PageData struct {
	Artists         []modules.Artist
	Q               string
	CreationDateMin int
	CreationDateMax int
	FirstAlbumMin   int
	FirstAlbumMax   int
	MembersMin      int
	MembersMax      int
	Locations       string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles("web/home.html"))
	pd := PageData{
		Artists:         Data.Artists,
		CreationDateMin: 1900,
	}
	if err := tmpl.Execute(w, pd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	query := r.URL.Query().Get("q")
	creationDateMin, _ := strconv.Atoi(r.URL.Query().Get("creationDateMin"))
	creationDateMax, _ := strconv.Atoi(r.URL.Query().Get("creationDateMax"))
	firstAlbumMin, _ := strconv.Atoi(r.URL.Query().Get("firstAlbumMin"))
	firstAlbumMax, _ := strconv.Atoi(r.URL.Query().Get("firstAlbumMax"))
	membersMin, _ := strconv.Atoi(r.URL.Query().Get("membersMin"))
	membersMax, _ := strconv.Atoi(r.URL.Query().Get("membersMax"))

	locationsParam := r.URL.Query().Get("locations")
	var locations []string
	if locationsParam != "" {
		for _, loc := range strings.Split(locationsParam, ",") {
			if t := strings.TrimSpace(loc); t != "" {
				locations = append(locations, t)
			}
		}
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
	pd := PageData{
		Artists:         results,
		Q:               query,
		CreationDateMin: creationDateMin,
		CreationDateMax: creationDateMax,
		FirstAlbumMin:   firstAlbumMin,
		FirstAlbumMax:   firstAlbumMax,
		MembersMin:      membersMin,
		MembersMax:      membersMax,
		Locations:       locationsParam,
	}

	tmpl := template.Must(template.ParseFiles("web/home.html"))
	if err := tmpl.Execute(w, pd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles("web/home.html"))
	pd := PageData{Artists: Data.Artists}
	if err := tmpl.Execute(w, pd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
