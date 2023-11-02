package geoweb

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"

	doparkscraper "github.com/dereulenspiegel/dopark-scraper"
	"github.com/gocolly/colly/v2"
)

const baseUrl = "https://geoweb1.digistadtdo.de/OWSServiceProxy/client/parken.jsp"

type Space struct {
	Name  string
	Total int
	Free  int
}

type Scraper struct {
	c   *colly.Collector
	log *slog.Logger
}

func NewScraper(log *slog.Logger) (*Scraper, error) {
	s := &Scraper{
		c:   colly.NewCollector(colly.AllowURLRevisit()),
		log: log,
	}
	return s, nil
}

func (s *Scraper) Scrape() (spaces []doparkscraper.Parking, err error) {
	spacesRegex := regexp.MustCompile(`^(\d+) Plätze von (\d+) frei$`)
	s.c.OnHTML("#infos > dl", func(e *colly.HTMLElement) {

		name := e.ChildText("dt > a")

		spacesString := e.ChildText(".plaetze")

		if spacesString == "keine Angaben" {
			s.log.Debug("no data for parking space", "spaceName", name)
			return
		}
		if spacesRegex.MatchString(spacesString) {
			m := spacesRegex.FindAllStringSubmatch(spacesString, -1)

			freeSpacesString := m[0][1]
			totalSpacesString := m[0][2]

			freeSpaces, err := strconv.Atoi(freeSpacesString)
			if err != nil {
				s.log.Error("failed to parse string for free spaces", "string", freeSpacesString, "err", err)
				return
			}
			totalSpaces, err := strconv.Atoi(totalSpacesString)
			if err != nil {
				s.log.Error("failed to parse string for total spaces", "string", totalSpacesString, "err", err)
				return
			}
			spaces = append(spaces, doparkscraper.Parking{
				Name:        name,
				FreeSpaces:  freeSpaces,
				TotalSpaces: totalSpaces,
			})
		} else {
			s.log.Warn("unparseable parking space status", "string", spacesString)
		}
	})
	s.c.OnError(func(r *colly.Response, cErr error) {
		s.log.Error("encountered error during visit: %s", cErr)
		err = fmt.Errorf("colly visit error: %w", cErr)
	})
	if err := s.c.Visit(baseUrl); err != nil {
		return nil, fmt.Errorf("failed visting the url %s: %s", baseUrl, err)
	}
	return spaces, nil
}
