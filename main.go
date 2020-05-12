package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/umahmood/haversine"
)

func main() {
	// Ensure the port has been set
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// crete a router using MUX and process the GET request
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", getSchools).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func GetSchools(w http.ResponseWriter, r *http.Request) {

	// Get the latitude URL paramater
	latString, ok := r.URL.Query()["lat"]
	if !ok || len(latString[0]) < 1 {
		respondWithError(w, http.StatusBadRequest, "Url Param 'lat' is missing")
		return
	}

	// Get the longitude URL paramater
	lonString, ok := r.URL.Query()["lon"]
	if !ok || len(lonString[0]) < 1 {
		respondWithError(w, http.StatusBadRequest, "Url Param 'lon' is missing")
		return
	}

	// Convert the lat & lon to strings
	latFloat, err := strconv.ParseFloat(latString[0], 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Please supply a valid lat paramater")
		return
	}

	lonFloat, err := strconv.ParseFloat(lonString[0], 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Please supply a valid lon paramater")
		return
	}

	currentLocation := haversine.Coord{Lat: latFloat, Lon: lonFloat}

	// Make Request to the code.org Schools API
	body, err := MakeRequest("https://code.org/schools.json")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to connect to Code.org")
		return
	}

	// parse the JSON response body
	var schools School
	json.Unmarshal(body, &schools)

	// Loop through the schools and calculate the distance from the supplied cordinates
	for i := 0; i < len(schools.Schools); i++ {
		schoolLocation := haversine.Coord{Lat: schools.Schools[i].Latitude, Lon: schools.Schools[i].Longitude}
		mi, km := haversine.Distance(currentLocation, schoolLocation)
		schools.Schools[i].Distance = mi
		schools.Schools[i].DistanceKM = km
	}

	// Sort the schools by their distance from the supplied cordinates
	sort.Slice(schools.Schools, func(i, j int) bool { return schools.Schools[i].Distance < schools.Schools[j].Distance })

	// Loop through the schools and get the nearest three for the response
	// todo: find a way to reuse the same struct instead of creating a duplicate
	var finalSchools []IsolatedSchools
	for i := 0; i < 3; i++ {
		finalSchools = append(finalSchools, schools.Schools[i])
	}
	// Log and return the result
	respondWithJSON(w, http.StatusOK, finalSchools)
}

// A struct used for parsing the http response
type School struct {
	Description string `json:"description"`
	Generated   string `json:"generated"`
	License     string `json:"license"`
	Schools     []struct {
		Name              string      `json:"name"`
		Website           string      `json:"website"`
		Levels            []string    `json:"levels"`
		Format            string      `json:"format"`
		FormatDescription string      `json:"format_description"`
		Gender            string      `json:"gender"`
		Description       string      `json:"description"`
		Languages         []string    `json:"languages"`
		MoneyNeeded       bool        `json:"money_needed"`
		OnlineOnly        bool        `json:"online_only"`
		NumberOfStudents  interface{} `json:"number_of_students"`
		ContactName       string      `json:"contact_name"`
		ContactNumber     string      `json:"contact_number"`
		ContactEmail      string      `json:"contact_email"`
		Latitude          float64     `json:"latitude"`
		Longitude         float64     `json:"longitude"`
		Street            string      `json:"street"`
		City              string      `json:"city"`
		State             string      `json:"state"`
		Zip               string      `json:"zip"`
		Published         int         `json:"published"`
		UpdatedAt         time.Time   `json:"updated_at"`
		Country           string      `json:"country"`
		Source            string      `json:"source"`
		Distance          float64     `json:"distance"`
		DistanceKM        float64     `json:"distanceKM"`
	} `json:"schools"`
}

// A struct used for soring the sorted school data
type IsolatedSchools struct {
	Name              string      `json:"name"`
	Website           string      `json:"website"`
	Levels            []string    `json:"levels"`
	Format            string      `json:"format"`
	FormatDescription string      `json:"format_description"`
	Gender            string      `json:"gender"`
	Description       string      `json:"description"`
	Languages         []string    `json:"languages"`
	MoneyNeeded       bool        `json:"money_needed"`
	OnlineOnly        bool        `json:"online_only"`
	NumberOfStudents  interface{} `json:"number_of_students"`
	ContactName       string      `json:"contact_name"`
	ContactNumber     string      `json:"contact_number"`
	ContactEmail      string      `json:"contact_email"`
	Latitude          float64     `json:"latitude"`
	Longitude         float64     `json:"longitude"`
	Street            string      `json:"street"`
	City              string      `json:"city"`
	State             string      `json:"state"`
	Zip               string      `json:"zip"`
	Published         int         `json:"published"`
	UpdatedAt         time.Time   `json:"updated_at"`
	Country           string      `json:"country"`
	Source            string      `json:"source"`
	Distance          float64     `json:"distance"`
	DistanceKM        float64     `json:"distanceKM"`
}

func MakeRequest(url string) ([]byte, error) {

	// Make the initila request to the Code.org API
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// Read the body of the callout response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Reaturn the response body
	return []byte(body), err
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	//encode payload to json
	response, _ := json.Marshal(payload)

	// set headers and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
