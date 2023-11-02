package doparkscraper

import (
	geom "github.com/twpayne/go-geom"
)

type Parking struct {
	Coordinates *geom.Point
	TotalSpaces int `db:"total_spaces"`
	FreeSpaces  int `db:"free_spaces"`
	Name        string
	Number      int
}
