package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"unsafe"

	_ "github.com/heroku/x/hmetrics/onload"
)

func TestGetSchoolsNoLat(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("lat", "")
	q.Add("lon", "-71.10869")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetSchools)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"error":"Url Param 'lat' is missing"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetSchoolsNoLon(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("lat", "42.959595")
	q.Add("lon", "")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetSchools)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"error":"Url Param 'lon' is missing"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetSchoolsResponseSize(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("lat", "42.959595")
	q.Add("lon", "-71.10869")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetSchools)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var testResponse testResponse
	body, err := ioutil.ReadAll(rr.Body)
	json.Unmarshal(body, &testResponse)

	assertEqual(t, unsafe.Sizeof(testResponse), 3, "The correct number of schools was not returned")

}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

type testResponse struct {
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
