package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	doparkscraper "github.com/dereulenspiegel/dopark-scraper"
	"github.com/dereulenspiegel/dopark-scraper/db"
	"github.com/dereulenspiegel/dopark-scraper/opendata"
	_ "github.com/lib/pq"
)

type scraper interface {
	Scrape() (spaces []doparkscraper.Parking, err error)
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log := slog.Default()

	ctx := context.Background()
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		dbUrl := os.Getenv("DOPARK_DB_URL")
		scrapeInterval := os.Getenv("DOPARK_INTERVAL")
		interval, err := time.ParseDuration(scrapeInterval)
		if err != nil {
			log.Error("invalid interval", "interval", scrapeInterval, "err", err)
			os.Exit(1)
		}
		pdb, err := sql.Open("postgres", dbUrl)
		if err != nil {
			log.Error("failed to connect to db", "err", err)
			os.Exit(1)
		}

		if err := db.RunMigrations(pdb); err != nil {
			log.Error("failed to run migrations", "err", err)
			os.Exit(1)
		}
		store, err := db.NewStore(pdb, log)
		if err != nil {
			log.Error("failed to create datastore", "err", err)
			os.Exit(1)
		}

		geoScraper := opendata.NewScraper(log)
		log.Info("Running scraper with interval", "interval", interval)
		ticker := time.NewTicker(interval)
		storeCtx, storeCancel := context.WithCancel(cancelCtx)
		defer storeCancel()
		for {
			select {
			case <-cancelCtx.Done():
				log.Info("scraping cancelled")
				return
			case <-ticker.C:
				log.Debug("starting scrape run")
				scrapeAndInsert(storeCtx, log, geoScraper, store)
			}
		}

	}()
	<-sigs
	log.Info("exiting")
}

func scrapeAndInsert(ctx context.Context, log *slog.Logger, scraper scraper, store *db.Store) {
	spaces, err := scraper.Scrape()
	if err != nil {
		log.Error("failed to scrape data", "error", err)
		return
	}
	for _, space := range spaces {
		if err := store.UpsertMetadata(ctx, &space); err != nil {
			log.Error("failed to insert metadata", "err", err)
		}
		log.Info("Inserting data", "id", space.Number)
		if err := store.InsertValues(ctx, &space); err != nil {
			log.Error("failed to insert values", "err", err)
		}
	}
}
