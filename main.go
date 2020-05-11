package main

import (
	"encoding/json"
	"fmt"
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

func getSchools(w http.ResponseWriter, r *http.Request) {

	// Get the latitude URL paramater
	latString, ok := r.URL.Query()["lat"]
	if !ok || len(latString[0]) < 1 {
		log.Println("Url Param 'lat' is missing")
		http.Error(w, "Bad Request: Url Param 'lat' is missing", 400)
		json.NewEncoder(w).Encode(http.Error)
		return
	}

	// Get the longitude URL paramater
	lonString, ok := r.URL.Query()["lon"]
	if !ok || len(lonString[0]) < 1 {
		http.Error(w, "Bad Request: Url Param 'lon' is missing", 400)
		json.NewEncoder(w).Encode(http.Error)
		return
	}

	// Convert the lat & lon to strings
	latFloat, err := strconv.ParseFloat(latString[0], 64)
	lonFloat, err := strconv.ParseFloat(lonString[0], 64)
	currentLocation := haversine.Coord{Lat: latFloat, Lon: lonFloat}

	// Make Request to the code.org Schools API
	body, err := MakeRequest("https://code.org/schools.json")
	if err != nil {
		http.Error(w, err.Error(), 400)
		fmt.Println(err)
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
		fmt.Println("User Type: " + schools.Schools[i].Name)
		fmt.Printf("%f Miles\n", schools.Schools[i].Distance)
	}

	// Take the list of nearby schools & serialize the response
	finalSchoolsAsString, err := json.Marshal(finalSchools)
	if err != nil {
		panic(err)
	}

	// Log and return the result
	fmt.Println(string(finalSchoolsAsString))
	json.NewEncoder(w).Encode(finalSchools)

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
