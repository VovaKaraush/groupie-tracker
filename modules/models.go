package modules

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

type Location struct {
	ID    int          `json:"id"`
	Towns []string     `json:"locations"`
	Dates ConcertDates `json:"dates"`
}

type Locations struct {
	Location Location
}

type ConcertDates struct {
	ID    int    `json:"id"`
	Dates string `json:"dates"`
}

type Relation struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}
