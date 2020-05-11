package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
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
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+port, router))

}

func homeLink(w http.ResponseWriter, r *http.Request) {
	body, err := MakeRequest("https://code.org/schools.json")
	if err != nil {
		fmt.Println(err)
	}

	latString, ok := r.URL.Query()["lat"]
	if !ok || len(latString[0]) < 1 {
		log.Println("Url Param 'lat' is missing")
		return
	}

	lonString, ok := r.URL.Query()["lon"]
	if !ok || len(lonString[0]) < 1 {
		log.Println("Url Param 'lat' is missing")
		return
	}

	latFloat, err := strconv.ParseFloat(latString[0], 64)
	lonFloat, err := strconv.ParseFloat(lonString[0], 64)

	var schools School

	json.Unmarshal(body, &schools)

	currentLocation := haversine.Coord{Lat: latFloat, Lon: lonFloat}

	for i := 0; i < len(schools.Schools); i++ {
		schoolLocation := haversine.Coord{Lat: schools.Schools[i].Latitude, Lon: schools.Schools[i].Longitude}
		mi, km := haversine.Distance(currentLocation, schoolLocation)
		schools.Schools[i].Distance = mi
	}

	sort.Slice(schools.Schools, func(i, j int) bool { return schools.Schools[i].Distance < schools.Schools[j].Distance })

	var finalSchools []IsolatedSchools

	for i := 0; i < 3; i++ {
		finalSchools = append(finalSchools, schools.Schools[i])
		fmt.Println("User Type: " + schools.Schools[i].Name)
		fmt.Printf("%f Miles\n", schools.Schools[i].Distance)
	}
	finalSchoolsAsString, err := json.Marshal(finalSchools)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(finalSchoolsAsString))
	json.NewEncoder(w).Encode(finalSchools)

}

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
		Distance          float64
	} `json:"schools"`
}

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
	Distance          float64
}

func MakeRequest(url string) ([]byte, error) {

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return []byte(body), err
}

func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}
