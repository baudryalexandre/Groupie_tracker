package main

import (
	"log"
	"net/http"

	groupie "groupie/utils"
)

func main() {
	fileServerCss := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServerCss))
	http.HandleFunc("/", groupie.MainPageHandler)
	http.HandleFunc("/band/", groupie.PageArtistHandler)
	http.HandleFunc("/search/", groupie.SearchHandler)
	http.HandleFunc("/filter", groupie.FilterHandler)

	// Adresse et port du serveur
	addr := ":8080"

	// Démarrage du serveur web.
	go func() {
		log.Println("Server started on http://localhost" + addr + "/groupie-tracker")
	}()

	http.ListenAndServe(addr, nil)
}
