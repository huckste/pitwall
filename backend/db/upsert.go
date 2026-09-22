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

func UpsertTeam(ctx context.Context, pool *pgxpool.Pool, t Team) error {

	_, err := pool.Exec(ctx, `
		INSERT INTO teams (ocb_id, series, name, short_name, full_name)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (ocb_id, series) DO UPDATE
		SET name = EXCLUDED.name, short_name = EXCLUDED.short_name, full_name = EXCLUDED.full_name
	`, t.OCBID, t.Series, t.Name, t.ShortName, t.FullName)

	return err
}
