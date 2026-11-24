package ingest

import (
	"context"
	"fmt"
	"log"
	"pitwall/backend/db"
	"pitwall/backend/models"
	"pitwall/backend/ocb"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SyncSeason(ctx context.Context, pool *pgxpool.Pool, series models.Series, year int) error {
	// --- seasons ---
	seasons, err := ocb.FetchSeasons(string(series))
	if err != nil {
		return err
	}

	for _, s := range seasons {
		if err := db.UpsertSeason(ctx, pool, models.MapSeason(s, series)); err != nil {
			return err
		}
	}
	log.Printf("upserted %d seasons for %s\n", len(seasons), series)

	var season *ocb.Season
	for _, s := range seasons {
		if s.Year == year {
			season = &s
			break
		}
	}

	if season == nil {
		return fmt.Errorf("no %d season found for %s", year, series)
	}

	// --- teams ---
	teams, err := ocb.FetchSeasonTeams(string(series), season.ID)
	if err != nil {
		return err
	}

	for _, t := range teams {
		if err := db.UpsertTeam(ctx, pool, models.MapSeasonTeam(t, series)); err != nil {
			return err
		}
	}
	log.Printf("upserted %d teams for %s %d\n", len(teams), series, year)

	// --- drivers ---
	drivers, err := ocb.FetchSeasonDrivers(string(series), season.ID)
	if err != nil {
		return err
	}

	for _, d := range drivers {
		if err := db.UpsertDriver(ctx, pool, models.MapSeasonDriver(d, series)); err != nil {
			return err
		}
	}
	log.Printf("upserted %d drivers for %s %d\n", len(drivers), series, year)

	// --- events ---
	events, err := ocb.FetchSeasonEvents(string(series), season.Year)
	if err != nil {
		return err
	}

	sessionCount := 0
	resultCount := 0

	for _, e := range events {
		log.Printf("syncing %s (%d sessions)...", e.Name, len(e.Schedule))

		ev, err := models.MapEvent(e, series, season.ID)
		if err != nil {
			return err
		}

		if err := db.UpsertEvent(ctx, pool, ev); err != nil {
			return err
		}

		// --- sessions ---
		for _, sess := range e.Schedule {
			session, err := models.MapSession(sess, series, e.ID)
			if err != nil {
				return err
			}

			if err := db.UpsertSession(ctx, pool, session); err != nil {
				return err
			}
			sessionCount++

			if sess.Status != "completed" {
				continue
			}

			has, err := db.SessionResultsFetched(ctx, pool, sess.ID)
			if err != nil {
				return err
			}

			if has {
				continue
			}

			time.Sleep(1100 * time.Millisecond)

			// --- results ---
			results, err := ocb.FetchSessionResults(string(series), e.ID, sess.ID)
			if err != nil {
				return err
			}

			for _, r := range results {

				if err := db.InsertMissingDriver(ctx, pool, models.MapResultDriver(r, series)); err != nil {
					log.Printf("skipping driver %s: %v", r.Driver.ID, err)
					continue
				}

				res, err := models.MapResult(r, series, sess.ID)
				if err != nil {
					return err
				}

				if err := db.UpsertResult(ctx, pool, res); err != nil {
					log.Printf("skipping result (session %s, driver %s): %v", sess.ID, r.Driver.ID, err)
					continue
				}
				resultCount++
			}

			if err := db.MarkResultsFetched(ctx, pool, sess.ID); err != nil {
				return err
			}
		}
	}

	log.Printf("upserted %d sessions for %s %d\n", sessionCount, series, year)
	log.Printf("upserted %d results for %s %d\n", resultCount, series, year)
	log.Printf("upserted %d events for %s %d\n", len(events), series, year)

	return nil
}
