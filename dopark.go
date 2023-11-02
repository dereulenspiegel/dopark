package doparkscraper

import (
	"time"

	geom "github.com/twpayne/go-geom"
)

type Parking struct {
	Coordinates *geom.Point
	TotalSpaces int `db:"total_spaces"`
	FreeSpaces  int `db:"free_spaces"`
	Name        string
	Number      int
	UpdatedAt   time.Time `db:"updated_at"`
}
