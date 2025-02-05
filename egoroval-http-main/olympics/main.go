//go:build !solution

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

var dataPath string = ""

type Ath struct {
	Athlete string `json:"athlete"`
	Age     int    `json:"age"`
	Country string `json:"country"`
	Year    int    `json:"year"`
	Date    string `json:"date"`
	Sport   string `json:"sport"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}

type AthI struct {
	Athlete       string      `json:"athlete"`
	Country       string      `json:"country"`
	Medals        MAY         `json:"medals"`
	MedalsByYears map[int]MAY `json:"medals_by_year"`
}

type MAY struct {
	Gold   int `json:"gold"`
	Silver int `json:"silver"`
	Bronze int `json:"bronze"`
	Total  int `json:"total"`
}

type TCIY struct {
	Country string `json:"country"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}

var aths []Ath

func parseJSON(w *http.ResponseWriter) {
	content, err := ioutil.ReadFile(dataPath)
	if err != nil {
		http.Error(*w, "", 400)
	}
	err = json.Unmarshal(content, &aths)
	if err != nil {
		panic(err)
	}
}

func searchByName(w *http.ResponseWriter, searchName string, sport string) AthI {
	var variants []Ath
	currentCountry := ""
	targetName := searchName
	i := 0
	for _, value := range aths {
		if value.Athlete == targetName {
			if sport == "" || value.Sport == sport {
				if i == 0 {
					currentCountry = value.Country
					variants = append(variants, value)
				} else {
					if value.Country == currentCountry {
						variants = append(variants, value)
					}
				}
				i++
			}
		}
	}
	if len(variants) == 0 {
		http.Error(*w, "", 404)
		(*w).WriteHeader(404)
		return AthI{}
	}
	var currentInfo AthI
	currentInfo.Athlete = targetName
	currentInfo.Country = currentCountry
	var medalsByYearMap = make(map[int]MAY)
	for _, value := range variants {
		currentInfo.Medals.Gold += value.Gold
		currentInfo.Medals.Silver += value.Silver
		currentInfo.Medals.Bronze += value.Bronze
		currentInfo.Medals.Total += value.Total
		medalsByYearMap[value.Year] = MAY{
			Gold:   medalsByYearMap[value.Year].Gold + value.Gold,
			Silver: medalsByYearMap[value.Year].Silver + value.Silver,
			Bronze: medalsByYearMap[value.Year].Bronze + value.Bronze,
			Total:  medalsByYearMap[value.Year].Total + value.Total,
		}
	}
	currentInfo.MedalsByYears = medalsByYearMap
	return currentInfo
}

func topInSports(w *http.ResponseWriter, currentSport string, limit int) []AthI {
	usedNames := make(map[string]bool)
	currentAthletes := make([]AthI, 0)
	for _, value := range aths {
		if value.Sport == currentSport && !usedNames[value.Athlete] {
			currentAthletes = append(currentAthletes, searchByName(w, value.Athlete, value.Sport))
			usedNames[value.Athlete] = true
		}
	}
	if len(currentAthletes) == 0 {
		http.Error(*w, "", 404)
		(*w).WriteHeader(404)
		return make([]AthI, 0)
	}
	sort.Slice(currentAthletes, func(i, j int) bool {
		if currentAthletes[i].Medals.Gold == currentAthletes[j].Medals.Gold {
			if currentAthletes[i].Medals.Silver == currentAthletes[j].Medals.Silver {
				if currentAthletes[i].Medals.Bronze == currentAthletes[j].Medals.Bronze {
					return currentAthletes[i].Athlete < currentAthletes[j].Athlete
				}
				return currentAthletes[i].Medals.Bronze > currentAthletes[j].Medals.Bronze
			}
			return currentAthletes[i].Medals.Silver > currentAthletes[j].Medals.Silver
		}
		return currentAthletes[i].Medals.Gold > currentAthletes[j].Medals.Gold
	})
	if len(currentAthletes) < limit {
		return currentAthletes
	}
	return currentAthletes[:limit]
}

func topCountriesInYear(w *http.ResponseWriter, year int, limit int) []TCIY {
	countriesMap := make(map[string]TCIY)
	for _, value := range aths {
		if value.Year == year {
			countriesMap[value.Country] = TCIY{
				Country: value.Country,
				Gold:    countriesMap[value.Country].Gold + value.Gold,
				Silver:  countriesMap[value.Country].Silver + value.Silver,
				Bronze:  countriesMap[value.Country].Bronze + value.Bronze,
				Total:   countriesMap[value.Country].Total + value.Total,
			}
		}
	}
	if len(countriesMap) == 0 {
		http.Error(*w, "", 404)
		(*w).WriteHeader(404)
		return make([]TCIY, 0)
	}
	sortedCountries := make([]TCIY, 0)
	countryUsed := make(map[string]bool)
	for _, value := range aths {
		if !countryUsed[value.Country] {
			sortedCountries = append(sortedCountries, countriesMap[value.Country])
			countryUsed[value.Country] = true
		}
	}
	sort.Slice(sortedCountries, func(i, j int) bool {
		if sortedCountries[i].Gold == sortedCountries[j].Gold {
			if sortedCountries[i].Silver == sortedCountries[j].Silver {
				if sortedCountries[i].Bronze == sortedCountries[j].Bronze {
					return sortedCountries[i].Country < sortedCountries[j].Country
				}
				return sortedCountries[i].Bronze > sortedCountries[j].Bronze
			}
			return sortedCountries[i].Silver > sortedCountries[j].Silver
		}
		return sortedCountries[i].Gold > sortedCountries[j].Gold
	})
	sortedNonNull := make([]TCIY, 0)
	for _, value := range sortedCountries {
		if value.Country != "" {
			sortedNonNull = append(sortedNonNull, value)
		}
	}
	if len(sortedNonNull) < limit {
		return sortedNonNull
	}
	return sortedNonNull[:limit]
}

func handler(w http.ResponseWriter, r *http.Request) {
	urlPath := "http://" + r.Host + r.URL.String()
	u, err := url.Parse(urlPath)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	parseJSON(&w)
	if len(q["name"]) > 0 {
		currentInfo := searchByName(&w, q["name"][0], "")
		if currentInfo.Athlete != "" {
			currentInfoJSON, _ := json.Marshal(currentInfo)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			_, _ = fmt.Fprintln(w, string(currentInfoJSON))
		}
	} else if len(q["sport"]) > 0 {
		limit := 3
		if len(q["limit"]) > 0 {
			limit, err = strconv.Atoi(q["limit"][0])
			if err != nil {
				http.Error(w, "", 400)
				w.WriteHeader(400)
			}
		}
		currentTop := topInSports(&w, q["sport"][0], limit)
		if len(currentTop) > 0 {
			currentTopJSON, _ := json.Marshal(currentTop)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			_, _ = fmt.Fprintln(w, string(currentTopJSON))
		}
	} else if len(q["year"]) > 0 {
		limit := 3
		if len(q["limit"]) > 0 {
			limit, err = strconv.Atoi(q["limit"][0])
			if err != nil {
				http.Error(w, "", 400)
				w.WriteHeader(400)
			}
		}
		year, _ := strconv.Atoi(q["year"][0])
		currentTop := topCountriesInYear(&w, year, limit)
		if len(currentTop) > 0 {
			currentTopJSON, _ := json.Marshal(currentTop)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			_, _ = fmt.Fprintln(w, string(currentTopJSON))
		}
	}
}

func main() {
	portPtr := flag.Int("port", rand.Intn(10000), "port string")
	dataPtr := flag.String("data", "", "data string")
	flag.Parse()
	portNumber := *portPtr
	dataPath = *dataPtr
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(portNumber), nil))
}
