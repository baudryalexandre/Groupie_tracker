package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ManageLocations(searchBoxValue string, artistsName []string) bool {
	/*
		Called for managing search and things that had matched

		searchBoxValue (string) is the user typed
		artists []string contains all artists member
	*/
	var doesArtistMatched bool = false
	GetArtistMatched(searchBoxValue, artistsName)
	return doesArtistMatched
}

func GetAllArtistName(artist []Artists) []string {
	/*
		Get an array of all artists name
	*/
	aResult := make([]string, 0)
	for _, a := range artist {
		aResult = append(aResult, a.Name)
	}
	return aResult
}

func GetArtistMatched(searchBoxValue string, arrayArtistsName []string) []string {
	/*
		Do a simple test with artist,

		Take args as:
			textbox-search value (string) that contains what user typed
			arrayArtistsName[]string contains all artists name
	*/
	for _, artistName := range arrayArtistsName {
		// Test si la saisie utilisateur (textbox-search value) est un duplitaca du nom d'artiste
		fmt.Println("artistName:'", artistName, "' est égale '", searchBoxValue, "'?")
		if artistName == searchBoxValue {
		}
	}
	return make([]string, 0)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	searchTerm := r.URL.Query().Get("search")

	// Convertit le terme de recherche en minuscules pour une recherche insensible à la casse
	searchTerm = strings.ToLower(searchTerm)

	groupID, found := GetGroupIDBySearchTerm(searchTerm)

	if !found {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Redirigez vers la page du groupe
	http.Redirect(w, r, "/band/?ID="+groupID, http.StatusFound)
}

func GetGroupIDBySearchTerm(searchTerm string) (string, bool) {
	// Interrogez votre API pour obtenir les correspondances entre les noms de groupe, d'artistes, de membres, de dates de création et de premiers albums
	url := "https://groupietrackers.herokuapp.com/api/artists" // Remplacez par l'URL correcte
	resp, err := http.Get(url)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()

	var artists []Artists // Assurez-vous que la structure Artists correspond à celle de votre API
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return "", false
	}

	// Convertissez le terme de recherche en minuscules pour une recherche insensible à la casse
	searchTermLower := strings.ToLower(searchTerm)

	// Parcourez la liste des artistes pour rechercher le nom d'artiste, le nom de membre, la date de création ou le premier album
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), searchTermLower) {
			return strconv.Itoa(artist.ID), true
		}

		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), searchTermLower) {
				return strconv.Itoa(artist.ID), true
			}
		}

		// Si le terme de recherche est une année (4 chiffres), recherchez dans CreationDate
		if len(searchTerm) == 4 && searchTerm == strconv.Itoa(artist.CreationDate) {
			return strconv.Itoa(artist.ID), true
		}

		// Si le terme de recherche est une date de création complète, recherchez dans CreationDate
		if isValidDate(searchTerm) && searchTerm == time.Unix(int64(artist.CreationDate), 0).Format("02-01-2006") {
			return strconv.Itoa(artist.ID), true
		}

		// Si le terme de recherche est une date de premier album au format "JJ-MM-AAAA", recherchez dans FirstAlbum
		if isValidDate(searchTerm) && searchTerm == artist.FirstAlbum {
			return strconv.Itoa(artist.ID), true
		}
	}

	// Si le terme de recherche n'est pas trouvé parmi les artistes, les membres, les dates de création ou les premiers albums, retournez false
	return "", false
}

// Fonction pour valider le format de la date "JJ-MM-AAAA"
func isValidDate(date string) bool {
	_, err := time.Parse("02-01-2006", date)
	return err == nil
}