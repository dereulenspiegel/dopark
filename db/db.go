package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

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
	name, coords
) VALUES (
	$1, $2
) ON CONFLICT(name) DO NOTHING`

var selectId = `SELECT number FROM spaces WHERE name=$1`

func (s *Store) UpsertMetadata(ctx context.Context, parking *doparkscraper.Parking) error {
	_, err := s.db.ExecContext(ctx, upsertMetadata, parking.Name, &ewkb.Point{Point: parking.Coordinates.SetSRID(4326)})
	if err != nil {
		return fmt.Errorf("failed to upsert metadata: %s", err)
	}
	var id int
	err = s.db.Get(&id, selectId, parking.Name)
	if err != nil {
		s.log.Error("failed to query id", "err", err)
		return fmt.Errorf("failed to query id: %w", err)
	}
	parking.Number = id
	return nil
}

var insertValues = `INSERT INTO park_values(
		spaces_id, 
		free, 
		total, 
		time,
		updated_at
	)
	VALUES(
		:number, 
		:free_spaces, 
		:total_spaces, 
		NOW(),
		:updated_at
	) `

func (s *Store) InsertValues(ctx context.Context, parking *doparkscraper.Parking) error {
	if parking.UpdatedAt.IsZero() {
		parking.UpdatedAt = time.Now()
	}
	_, err := s.db.NamedExecContext(ctx, insertValues, parking)
	if err != nil {
		return fmt.Errorf("failed to insert values: %s", err)
	}
	return nil
}
