package doparkscraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	geom "github.com/twpayne/go-geom"
)

const prodEndpoint string = "https://www.dopark.de/phbelegung/PHBelegung.json"

type Properties struct {
	Name               string `json:"name"`
	TotalParkingSpaces int    `json:"plaetze"`
	FreeParkingSpaces  int    `json:"frei"`
	Phnummer           int    `json:"phnummer"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Feature struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
	Geometry   Geometry   `json:"geometry"`
}

type Collection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Parking struct {
	Lat         float64
	Lon         float64
	Coordinates geom.Point
	TotalSpaces int `db:"total_spaces"`
	FreeSpaces  int `db:"free_spaces"`
	Name        string
	Number      int
}

func Scrape() (spaces []Parking, err error) {
	resp, err := http.Get(prodEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %s", prodEndpoint, err)
	}
	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024*2))
	if err != nil {
		return nil, fmt.Errorf("failed to read data from %s: %s", prodEndpoint, err)
	}
	coll := Collection{}
	if err := json.Unmarshal(bodyBytes[3:], &coll); err != nil {
		return nil, fmt.Errorf("failed to parse response from %s: %s", prodEndpoint, err)
	}
	for _, feat := range coll.Features {
		coordinates := geom.NewPoint(geom.XY).MustSetCoords([]float64{feat.Geometry.Coordinates[0], feat.Geometry.Coordinates[1]}).SetSRID(25832)
		spaces = append(spaces, Parking{
			Lat:         feat.Geometry.Coordinates[1], // TODO check order of coordinates
			Lon:         feat.Geometry.Coordinates[0],
			TotalSpaces: feat.Properties.TotalParkingSpaces,
			FreeSpaces:  feat.Properties.FreeParkingSpaces,
			Name:        feat.Properties.Name,
			Number:      feat.Properties.Phnummer,
			Coordinates: *coordinates,
		})
	}
	return
}
