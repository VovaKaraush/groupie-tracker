package modules

// GroupieData regroupe toutes les données fetchées de l'API
type GroupieData struct {
	Artists   []Artist
	Locations []Location
	Dates     []DateInfo
	Relations []Relation
}

type Artist struct {
	ID               int      `json:"id"`
	Image            string   `json:"image"`
	Name             string   `json:"name"`
	Members          []string `json:"members"`
	CreationDate     int      `json:"creationDate"`
	FirstAlbum       string   `json:"firstAlbum"`
	ConcertLocations string   `json:"locations"`
	ConcertDates     string   `json:"concertDates"`
	Relations        string   `json:"relations"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type LocationsResponse struct {
	Index []Location `json:"index"`
}

type RelationResponse struct {
	Index []Relation `json:"index"`
}

type DateInfo struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type DatesResponse struct {
	Index []DateInfo `json:"index"`
}

type Relation struct {
	ID             int                    `json:"id"`
	DatesLocations map[string]interface{} `json:"datesLocations"`
}
