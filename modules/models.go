package main

// Models used by the API
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    []string `json:"locations"`
	ConcertDates []string `json:"concertDates"`
	Relations    int      `json:"relations"`
}

type Relation struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}
