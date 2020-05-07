package usgs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

type Response struct {
	Type     string     `json:"type"`
	Metadata Metadata   `json:"metadata"`
	Features []Features `json:"features"`
	Bbox     []float64  `json:"bbox"`
}
type Metadata struct {
	Generated int
	Url       string
	Title     string
	Status    int
	Api       string
	Count     int
}
type Features struct {
	Type       string
	Properties Properties
	Geometry   Geometry
	Id         string
}
type Properties struct {
	Mag     float64
	Place   string
	Time    int
	Updated int
	Tz      int
	Url     string
	Detail  string
	Felt    string
	Cdi     string
	Mmi     string
	Alert   string
	Status  string
	Tsunami int
	Sig     int
	Net     string
	Code    string
	Ids     string
	Sources string
	Types   string
	Nst     int
	Dmin    float32
	Rms     float32
	Gap     int
	MagType string
	Type    string
	Title   string
}

type Geometry struct {
	Type        string
	Coordinates []float32
}

// URLLastHourQuakes this is the url to obtain the quakes occurred in the last hour from the USGS database
const URLLastHourQuakes = "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_hour.geojson"

var myClient = &http.Client{Timeout: 10 * time.Second}

// GetLastHourQuakes Method to retrieve from a url(USGS) the list of the quakes in the last hour
func GetLastHourQuakes() {

	r, err := myClient.Get(URLLastHourQuakes)

	if err != nil {
		fmt.Println("Error while loading the file")
		fmt.Println(err.Error())
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	var quakesResponse Response
	json.Unmarshal(body, &quakesResponse)
	fmt.Println(reflect.TypeOf(quakesResponse))
	fmt.Println("------------------------------")
	fmt.Println(quakesResponse.Type)
	fmt.Println("------------------------------")
	fmt.Println(quakesResponse.Metadata.Url)
	fmt.Println("------------------------------")
	fmt.Println(quakesResponse.Metadata.Url)
	fmt.Println("------------------------------")
	fmt.Println(len(quakesResponse.Features))
	fmt.Println("------------------------------")

	for i := range quakesResponse.Features {
		fmt.Println(i)
		fmt.Println(quakesResponse.Features[i].Properties.Title)
		fmt.Println(reflect.TypeOf(quakesResponse.Features[i].Properties.Title))
	}

}
