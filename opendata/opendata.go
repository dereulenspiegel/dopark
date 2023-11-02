package opendata

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	doparkscraper "github.com/dereulenspiegel/dopark-scraper"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
)

const OpendataUrl = `https://geoweb1.digistadtdo.de/doris_gdi/geoserver/parken/ogc/features/collections/pls/items?f=application%2Fgeo%2Bjson`
const maxResponseSize = 5 * 1024 * 1024

type Scraper struct {
	log *slog.Logger
}

func NewScraper(log *slog.Logger) *Scraper {
	return &Scraper{
		log: log,
	}
}

func (s *Scraper) Scrape() (spaces []doparkscraper.Parking, err error) {
	resp, err := http.Get(OpendataUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to query url %s: %w", OpendataUrl, err)
	}
	defer resp.Body.Close()
	return scrape(s.log, io.LimitReader(resp.Body, maxResponseSize))
}

func scrape(log *slog.Logger, input io.Reader) (spaces []doparkscraper.Parking, err error) {
	var collection geojson.FeatureCollection
	if err := json.NewDecoder(input).Decode(&collection); err != nil {
		return nil, fmt.Errorf("failed to parse feature collection: %w", err)
	}
	for _, feature := range collection.Features {
		coords := feature.Geometry
		name := feature.Properties["name"]
		logger := log.With("name", name)
		totalSpacesRaw, ok := feature.Properties["cap"].(float64)
		if !ok {
			logger.Warn("invalid value type for total spaces", "value", feature.Properties["cap"])
			continue
		}
		freeSpacesRaw, ok := feature.Properties["frei"].(float64)
		if !ok {
			logger.Warn("invalid value type for free spaces", "value", feature.Properties["frei"])
			continue
		}
		updatesAt, err := time.Parse(time.DateTime, feature.Properties["stand"].(string))
		if err != nil {
			logger.Warn("invalid update date time", "dateTimeString", feature.Properties["stand"])
		}

		spaces = append(spaces, doparkscraper.Parking{
			Coordinates: geom.NewPointFlat(coords.Layout(), coords.FlatCoords()),
			TotalSpaces: int(totalSpacesRaw),
			FreeSpaces:  int(freeSpacesRaw),
			Name:        name.(string),
			UpdatedAt:   updatesAt,
		})
	}
	return spaces, nil
}
