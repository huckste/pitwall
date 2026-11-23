package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Team struct {
	OCBID     string
	Series    string
	Name      string
	ShortName string
	FullName  string
}

type Driver struct {
	OCBID       string
	Series      string
	FirstName   string
	LastName    string
	Number      *int
	Nationality *string
	TeamOCBID   *string
}

func UpsertTeam(ctx context.Context, pool *pgxpool.Pool, t Team) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO teams (ocb_id, series, name, short_name, full_name)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (ocb_id, series) DO UPDATE
		SET name = EXCLUDED.name, 
			short_name = EXCLUDED.short_name, 
			full_name = EXCLUDED.full_name
	`, t.OCBID, t.Series, t.Name, t.ShortName, t.FullName)

	return err
}

func UpsertDriver(ctx context.Context, pool *pgxpool.Pool, d Driver) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO drivers (ocb_id, series, first_name, last_name, number, nationality, team_id)
		VALUES ($1, $2, $3, $4, $5, $6,
			(SELECT id FROM teams WHERE ocb_id = $7 AND series = $2))
		ON CONFLICT (ocb_id, series) DO UPDATE
		SET first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			number = EXCLUDED.number,
			nationality = EXCLUDED.nationality,
			team_id = EXCLUDED.team_id
		`, d.OCBID, d.Series, d.FirstName, d.LastName, d.Number, d.Nationality, d.TeamOCBID)

	return err
}
