package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

const URL = "https://groupietrackers.herokuapp.com/api"

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/groupie-tracker" {
		// http.Redirect(w,r,"/404", 404)
		// http.Error(w, "Page inexistante.", 404)
		NotFoundHandler(w, r)
		return
	}
	SPEED(w,r)
	tmpl := template.Must(template.ParseFiles("index.html"))
	var data JSON
	resp, err := http.Get(URL)
	if err != nil {
		panic("erreur 500.css")
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&data)
	url := data.Artists
	Art, err := http.Get(url)
	if err != nil {
		panic("Une erreur est survenue:,JERem:")
	}
	defer Art.Body.Close()
	var Artists []Artists
	err = json.NewDecoder(Art.Body).Decode(&Artists)
	err = tmpl.Execute(w, Artists)

	if r.Method == "POST" {
		fmt.Println("Bouton pressé")
	}
	fmt.Println("alex-dieu:", r.Method)
	if err != nil {
		panic("Une erreur est survenue,Jer:" + err.Error())
	}
}

func PageArtistHandler(w http.ResponseWriter, r *http.Request) {
	artistID := r.FormValue("ID")
	fmt.Println("artistsID:", artistID)
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.ParseFiles("templates/band.html"))
		var data JSON
		resp, err := http.Get(URL)
		if err != nil {
			ServerErrorHandler(w, r)
			return
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&data)
		urlArtist := data.Artists + "/" + artistID
		Art, err := http.Get(urlArtist)
		if err != nil {
			ServerErrorHandler(w, r)
			return
		}
		defer Art.Body.Close()
		var Locations Locations
		urlLocations := data.Locations + "/" + artistID
		Loc, err := http.Get(urlLocations)
		if err != nil {
			ServerErrorHandler(w, r)
			return
		}
		err = json.NewDecoder(Loc.Body).Decode(&Locations)
		if err != nil {
			ServerErrorHandler(w, r)
			return
		}
		defer Loc.Body.Close()
		var artist Artist
		err = json.NewDecoder(Art.Body).Decode(&artist)
		urlRelation := "https://groupietrackers.herokuapp.com/api/relation/" + artistID
		Rel, err := http.Get(urlRelation)
		if err != nil {
			ServerErrorHandler(w, r)
			return
		} else {
			fmt.Println("Tentative d'ouverture:", urlRelation)
		}
		var relation Relation
		err = json.NewDecoder(Rel.Body).Decode(&relation)
		if err != nil {
			ServerErrorHandler(w, r)
			return
		}
		defer Rel.Body.Close()

		var outputBand OutputSpecificBand

		outputBand.ID = artist.ID
		outputBand.Image = artist.Image
		outputBand.FirstAlbum = artist.FirstAlbum
		outputBand.Name = artist.Name
		outputBand.Members = artist.Members
		outputBand.CreationDate = artist.CreationDate
		outputBand.Locations = Format_Locations_From_Array(Locations.Lcs)
		outputBand.RelationA = Format_Date(relation.DatesLocations)

		err = tmpl.Execute(w, outputBand)
		if err != nil {
			http.Redirect(w, r, "/500", http.StatusInternalServerError)
			return
		}
	case "POST":
		// do nothing
	default:
		http.Redirect(w, r, "400", 400)
	}
}

func Format_Date(m map[string][]string) []string {
	/*
		Utilisé pour formater vers une chaîne les concerts avec les dates
	*/
	a := make([]string, 0)
	var line string
	for k, v := range m {
		line = "\n" + Format_Location_From_String(k) + " : "
		for i, d := range v {
			line += d
			if i < len(v)-1 {
				line += "  " // Ajoutez un espace entre les dates
			}
		}
		a = append(a, line)
	}
	return a
}

func CinqCentHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	tmpl := template.Must(template.ParseFiles("templates/500.html"))
	tmpl.Execute(w, nil)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/404.html")
	if err != nil {
		panic("Niet.")
	}
	w.WriteHeader(http.StatusNotFound)
	tmpl.Execute(w, http.StatusNotFound)
}

func ServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		panic("Niet.")
	}
	w.WriteHeader(http.StatusInternalServerError)
	tmpl.Execute(w, nil)
}

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		tmpl := template.Must(template.ParseFiles("templates/filter.html"))
		var data JSON
		resp, err := http.Get(URL)
		if err != nil {
			ServerErrorHandler(w, r)
			return
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&data)
		urlArtist := data.Artists //+ "/" + artistID
		Art, err := http.Get(urlArtist)
		if err != nil {
			ServerErrorHandler(w, r)
			return
		}
		// defer Art.Body.Close()
		// var Locations Locations
		// urlLocations := data.Locations + "/" + artistID
		// Loc, err := http.Get(urlLocations)
		// if err != nil {
		// 	ServerErrorHandler(w, r)
		// 	return
		// }
		// err = json.NewDecoder(Loc.Body).Decode(&Locations)

		// defer Loc.Body.Close()
		var artists []Artists
		err = json.NewDecoder(Art.Body).Decode(&artists)
		// if err != nil {
		// 	fmt.Println("Mince une erreur est survenue:", err)
		// }
		// urlRelation := "https://groupietrackers.herokuapp.com/api/relation" //+ artistID
		// Rel, err := http.Get(urlRelation)
		// if err != nil {
		// 	ServerErrorHandler(w, r)
		// 	return
		// } else {
		// 	fmt.Println("Tentative d'ouverture:", urlRelation)
		// }
		// var relation Relation
		// err = json.NewDecoder(Rel.Body).Decode(&relation)
		// defer Rel.Body.Close()
		var outputBand OutputBand
		var arTists []Artists

		// fmt.Println("artists:", artists)
		// fmt.Println("r.Method:", r.Method)
		if r.Method == "POST" {
			r.ParseForm()
			mCheckBox := GetCheckBoxValue(r)
			startCreationYear := r.FormValue("startCreationYearRange")
			endCreationYear := r.FormValue("endCreationYearRange")
			startFirstAlbumYear := r.FormValue("startFirstAlbumYearRange")
			endFirstAlbumYear := r.FormValue("endFirstAlbumYearRange")
			mLocations := GetLocations(r)

			arTists = SortBand(artists, mCheckBox, mLocations, Atoi(startFirstAlbumYear), Atoi(endFirstAlbumYear), Atoi(startCreationYear), Atoi(endCreationYear), w, r)
		}

		outputBand.Artists = arTists

		err = tmpl.Execute(w, outputBand)
	case "GET":
		// do nothing
	default:
		http.Redirect(w, r, "400", 400)
	}
}

func SPEED(w http.ResponseWriter, r *http.Request) {
	index, _ := os.ReadFile("index.html")
	var nodata JSON
	resp, err := http.Get(URL)
	if err != nil {
		panic("erreur 500.css")
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&nodata)
	url := nodata.Artists
	Art, err := http.Get(url)
	if err != nil {
		panic("Une erreur est survenue:,JERem:")
	}
	defer Art.Body.Close()
	var Artists []Artists
	err = json.NewDecoder(Art.Body).Decode(&Artists)
	if err != nil {
		panic("Une erreur est survenue,Jer:" + err.Error())
	}

	// var Némo [][]string
	var tempArr []string
	var tempMem []string
	for _, artists := range Artists {
		// fmt.Println(artists.Name)
		var artistsName string = artists.Name
		artistsName = artists.Name + "- Artist"
		tempArr = append(tempArr, artistsName)
		for _, members := range artists.Members {
			artistsName = members + "-" + artists.Name
			tempMem = append(tempMem, artistsName)

		}
		// Némo = append(Némo, artists.Members)

	}
	html := string(index)
	optionsArtiste := ""
	optionsMembres := ""
	for _, item := range tempArr {
		optionsArtiste += fmt.Sprintf("<option value='%v'>%v</option>\n", item, item)
	}
	for _, item := range tempMem {
		optionsArtiste += fmt.Sprintf("<option value='%v'>%v</option>\n", item, item)
	}

	// for _, item := range Némo{
	// 	for _, item2 := range item{
	// 		optionsMembres += fmt.Sprintf("<option value='%v'>%v - Membres</option>\n", item2, item2)
	// 	}
	// }
	html = strings.Replace(html, "<!-- Groupe -->", optionsArtiste, 1)
	html = strings.Replace(html, "<!-- Membres -->", optionsMembres, 1)
	// Write the modified HTML to the response
	w.Header().Set("Content-Type", "text/html")
	os.WriteFile("index.html", []byte(html), 0o666)
}