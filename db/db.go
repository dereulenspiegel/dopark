package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	doparkscraper "github.com/dereulenspiegel/dopark-scraper"
	"github.com/jmoiron/sqlx"
	"github.com/twpayne/go-geom/encoding/ewkb"
)

type Store struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewStore(db *sql.DB, log *slog.Logger) (*Store, error) {
	return &Store{
		db:  sqlx.NewDb(db, "postgres"),
		log: log,
	}, nil
}

var upsertMetadata = `INSERT INTO spaces(
	name, coords, number
) VALUES (
	$1, $2, $3
) ON CONFLICT(name) DO NOTHING`

func (s *Store) UpsertMetadata(ctx context.Context, parking doparkscraper.Parking) error {
	_, err := s.db.ExecContext(ctx, upsertMetadata, parking.Name, &ewkb.Point{Point: &parking.Coordinates}, parking.Number)
	if err != nil {
		return fmt.Errorf("failed to upsert metadata: %s", err)
	}
	return nil
}

var insertValues = `INSERT INTO park_values(
		spaces_id, 
		free, 
		total, 
		time
	)
	VALUES(
		:number, 
		:free_spaces, 
		:total_spaces, 
		NOW()
	) `

func (s *Store) InsertValues(ctx context.Context, parking doparkscraper.Parking) error {
	_, err := s.db.NamedExecContext(ctx, insertValues, parking)
	if err != nil {
		return fmt.Errorf("failed to insert values: %s", err)
	}
	return nil
}
