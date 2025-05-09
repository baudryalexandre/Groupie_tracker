package utils

import (
	"fmt"
	"net/http"
	"regexp"
)

func SortCreationDate(arrayAllArtists []Artist, from_date, to_date int) []Artist {
	/*
		Return all bands that had started after from_date

	*/
	arrayArtists := make([]Artist, 0)
	for _, artist := range arrayAllArtists {
		if artist.CreationDate > from_date && artist.CreationDate < to_date {
			arrayArtists = append(arrayArtists, artist)
		}
	}
	return arrayArtists
}

func GetCheckBoxValue(r *http.Request) []bool {
	/*
		Called for getting checkbox value into an array,
		this array is returned.
		Theses checkbox elements are accessed by an request object
	*/
	a := make([]bool, 0)
	fmt.Println("firstValue:", r.FormValue("checkboxOneMember"))
	a = append(a, (string(r.FormValue("checkboxOneMember")) == "on"))
	a = append(a, (string(r.FormValue("checkboxTwoMembers")) == "on"))
	a = append(a, (string(r.FormValue("checkboxThreeMembers")) == "on"))
	a = append(a, (string(r.FormValue("checkboxFourMembers")) == "on"))
	a = append(a, (string(r.FormValue("checkboxFiveMembers")) == "on"))
	a = append(a, (string(r.FormValue("checkboxSixMembers")) == "on"))
	a = append(a, (string(r.FormValue("checkboxSevenMembers")) == "on"))
	a = append(a, (string(r.FormValue("checkboxEightMembers")) == "on"))
	return a
}

func SortBand(aArtists []Artists, mCheckBox []bool, mLocations map[string]string, startFirstAlbumYear, endFirstAlbumYear, startCreationDateYear, endCreationDateYear int, w http.ResponseWriter, r *http.Request) []Artists {
	/*
		Called for getting all artists sorted.
		All artists that have criteria that matched with those passed as arguments are appended in the array
	*/
	var arrayArtists []Artists
	isLocationsMatched := false
	arrayMembersValuesPossibiles := GetAllMembersValuePossibles(mCheckBox)
	arrayLocationsModel := GetLocationsAskedByUser(mLocations)
	for _, artist := range aArtists {
		doesFirstAlbumYearMatched := GetIsDateValid(artist.FirstAlbum, startFirstAlbumYear, endFirstAlbumYear)
		doesCreationYearMatched := GetIsDateValid(artist.CreationDate, startCreationDateYear, endCreationDateYear)

		isLocationsMatched = ManageLocation(artist.Locations, arrayLocationsModel, w, r)
		isEnoughMembers := GetMembersMatched(arrayMembersValuesPossibiles, len(artist.Members))

		if len(arrayLocationsModel) == 0 {
			isLocationsMatched = true
		}

		if doesFirstAlbumYearMatched && doesCreationYearMatched && isEnoughMembers && isLocationsMatched {
			arrayArtists = append(arrayArtists, artist)
		}
	}
	return arrayArtists
}

func GetMembersMatched(aCheckBox []int, countMembers int) bool {
	/*
		Called for getting if the number of members correspond to one of theses asked by the user
	*/
	for _, year := range aCheckBox {
		if year == countMembers {
			return true
		}
	}
	return false
}

func GetAllMembersValuePossibles(mCheckBoxIndexMembers []bool) []int {
	/*
		Called for getting all checkbox values from an array passed as argument into array of int
	*/
	a := make([]int, 0)
	var index int = 0
	for _, b := range mCheckBoxIndexMembers {
		if b {
			a = append(a, (index + 1))
		}
		index++
	}
	return a
}

func GetIsDateValid(year interface{}, startYear, endYear int) bool {
	switch year := year.(type) {
	case string:
		yearInt := GetYear(year)
		return yearInt >= startYear && yearInt <= endYear
	case int:
		return year >= startYear && year <= endYear
	default:
		return false
	}
}

func GetYear(year string) int {
	reg := regexp.MustCompile(`(\d{2})-(\d{2})-(?P<year>\d{4})`)
	a := reg.FindStringSubmatch(year)
	return Atoi(a[len(a)-1])
}

func Atoi(s string) int {
	var result int = 0
	var factor int = 1
	for i := len(s) - 1; i >= 0; i-- {
		digit := int(s[i] - '0') // Convertir le caractère en chiffre
		result += digit * factor
		factor *= 10
	}
	return result
}
