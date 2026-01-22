package handlers

import (
	"encoding/json"
	"groupie-tracker/filters"
	"groupie-tracker/modules"
	"net/http"
	"strconv"
	"text/template"
)

var Data *modules.GroupieData

func SortData(g *modules.GroupieData) {
	Data = g
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/home.html"))
	tmpl.Execute(w, nil)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query().Get("q")
	creationDateMin, _ := strconv.Atoi(r.URL.Query().Get("creationDateMin"))
	creationDateMax, _ := strconv.Atoi(r.URL.Query().Get("creationDateMax"))
	firstAlbumMin, _ := strconv.Atoi(r.URL.Query().Get("firstAlbumMin"))
	firstAlbumMax, _ := strconv.Atoi(r.URL.Query().Get("firstAlbumMax"))
	membersMin, _ := strconv.Atoi(r.URL.Query().Get("membersMin"))
	membersMax, _ := strconv.Atoi(r.URL.Query().Get("membersMax"))
	locations := r.URL.Query()["locations"]

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

	json.NewEncoder(w).Encode(results)
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/home.html"))
	tmpl.Execute(w, Data.Artists)
}
