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

var Data *modules.GroupieData

func SortData(g *modules.GroupieData) { Data = g }

// --------- SSR TEMPLATE (chargé 1 fois) ---------

type HomePageData struct {
	// Filters sert à ré-afficher les valeurs saisies dans le <form>
	// (on garde les valeurs brutes en string pour éviter d'afficher "0" dans les champs).
	Filters      FiltersView
	AllLocations []string
	Artists      []modules.Artist
}

type FiltersView struct {
	Q string

	CreationDateMin string
	CreationDateMax string
	FirstAlbumMin   string
	FirstAlbumMax   string
	MembersMin      string
	MembersMax      string
	Locations       string // champ texte: "Paris, Lyon"
}

var homeTmpl = template.Must(template.New("home.html").Funcs(template.FuncMap{
	"intOrEmpty": func(n int) string {
		if n == 0 {
			return ""
		}
		return strconv.Itoa(n)
	},
	"contains": func(list []string, item string) bool {
		for _, v := range list {
			if v == item {
				return true
			}
		}
		return false
	},
}).ParseFiles("web/home.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if Data == nil {
		http.Error(w, "Data not loaded yet", http.StatusInternalServerError)
		return
	}

	// Lire les params (GET)
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	creationDateMinRaw := strings.TrimSpace(r.URL.Query().Get("creationDateMin"))
	creationDateMaxRaw := strings.TrimSpace(r.URL.Query().Get("creationDateMax"))
	firstAlbumMinRaw := strings.TrimSpace(r.URL.Query().Get("firstAlbumMin"))
	firstAlbumMaxRaw := strings.TrimSpace(r.URL.Query().Get("firstAlbumMax"))
	membersMinRaw := strings.TrimSpace(r.URL.Query().Get("membersMin"))
	membersMaxRaw := strings.TrimSpace(r.URL.Query().Get("membersMax"))
	locationsRaw := strings.TrimSpace(r.URL.Query().Get("locations"))

	creationDateMin, _ := strconv.Atoi(creationDateMinRaw)
	creationDateMax, _ := strconv.Atoi(creationDateMaxRaw)
	firstAlbumMin, _ := strconv.Atoi(firstAlbumMinRaw)
	firstAlbumMax, _ := strconv.Atoi(firstAlbumMaxRaw)
	membersMin, _ := strconv.Atoi(membersMinRaw)
	membersMax, _ := strconv.Atoi(membersMaxRaw)

	// home.html a un champ texte "locations" (villes séparées par des virgules)
	locations := splitLocations(locationsRaw)

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

	// Si aucun filtre (et pas de query), on affiche tout
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
			Locations:       locationsRaw,
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


func SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query().Get("q")
	creationDateMin, _ := strconv.Atoi(r.URL.Query().Get("creationDateMin"))
	creationDateMax, _ := strconv.Atoi(r.URL.Query().Get("creationDateMax"))
	firstAlbumMin, _ := strconv.Atoi(r.URL.Query().Get("firstAlbumMin"))
	firstAlbumMax, _ := strconv.Atoi(r.URL.Query().Get("firstAlbumMax"))
	membersMin, _ := strconv.Atoi(r.URL.Query().Get("membersMin"))
	membersMax, _ := strconv.Atoi(r.URL.Query().Get("membersMax"))
	// Supporte à la fois :
	// - un champ texte "locations=Paris, Lyon"
	// - ou des checkboxes multiples "locations=Paris&locations=Lyon"
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

// Pour éviter de casser le code si ArtistHandler est encore routé quelque part :
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	if Data == nil {
		http.Error(w, "Data not loaded yet", http.StatusInternalServerError)
		return
	}
	page := HomePageData{
		Filters:      FiltersView{},
		AllLocations: uniqueLocations(Data),
		Artists:      Data.Artists,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = homeTmpl.Execute(w, page)
}

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