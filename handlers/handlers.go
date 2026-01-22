package handlers

import (
	"groupie-tracker/modules"
	"net/http"
	"text/template"
)

var Artists *modules.Artist

func SortData(g *modules.Artist) {
	Artists = g
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/home.html"))
	tmpl.Execute(w, Artists)
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/home.html"))
	tmpl.Execute(w, Artists)
}
