package main

import (
	"context"
	"log"
	"pitwall/backend/db"
	"pitwall/backend/models"
	"pitwall/backend/ocb"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	ctx := context.Background()

	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatal("failed to connect to db: ", err)
	}

	defer pool.Close()
	log.Printf("backfill connected to db, ready to go")

	seasons, err := ocb.FetchSeasons("formula1")
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range seasons {
		if err := db.UpsertSeason(ctx, pool, models.MapSeason(s, models.Formula1)); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("upserted %d seasons for formula1\n", len(seasons))

	var season2025 *ocb.Season

	for _, s := range seasons {
		if s.Year == 2025 {
			season2025 = &s
			break
		}
	}

	if season2025 == nil {
		log.Fatal("no 2025 season found for formula1")
	}

	events, err := ocb.FetchSeasonEvents("formula1", season2025.Year)
	if err != nil {
		log.Fatal(err)
	}

	sessionCount := 0
	resultCount := 0

	for _, e := range events {
		log.Printf("fetching results for %s (%d sessions)...", e.Name, len(e.Schedule))

		ev, err := models.MapEvent(e, models.Formula1, season2025.ID)
		if err != nil {
			log.Fatal(err)
		}

		if err := db.UpsertEvent(ctx, pool, ev); err != nil {
			log.Fatal(err)
		}

		for _, sess := range e.Schedule {
			session, err := models.MapSession(sess, models.Formula1, e.ID)
			if err != nil {
				log.Fatal(err)
			}

			if err := db.UpsertSession(ctx, pool, session); err != nil {
				log.Fatal(err)
			}
			sessionCount++

			if sess.Status != "completed" {
				continue
			}

			time.Sleep(1100 * time.Millisecond)

			results, err := ocb.FetchSessionResults("formula1", e.ID, sess.ID)
			if err != nil {
				log.Fatal(err)
			}

			for _, r := range results {

				if err := db.InsertMissingDriver(ctx, pool, models.MapResultDriver(r, models.Formula1)); err != nil {
					log.Printf("skipping driver %s: %v", r.Driver.ID, err)
					continue
				}

				res, err := models.MapResult(r, models.Formula1, sess.ID)
				if err != nil {
					log.Fatal(err)
				}

				if err := db.UpsertResult(ctx, pool, res); err != nil {
					log.Printf("skipping result (session %s, driver %s): %v", sess.ID, r.Driver.ID, err)
					continue
				}
				resultCount++
			}
		}
	}

	log.Printf("upserted %d sessions for formula1 2025\n", sessionCount)
	log.Printf("upserted %d results for formula1 2025\n", resultCount)
	log.Printf("upserted %d events for formula1 2025\n", len(events))

	teams, err := ocb.FetchSeasonTeams("formula1", season2025.ID)
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range teams {
		if err := db.UpsertTeam(ctx, pool, models.MapSeasonTeam(t, models.Formula1)); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("upserted %d teams for formula1 2025\n", len(teams))

	drivers, err := ocb.FetchSeasonDrivers("formula1", season2025.ID)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range drivers {
		if err := db.UpsertDriver(ctx, pool, models.MapSeasonDriver(d, models.Formula1)); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("upserted %d drivers for formula1 2025\n", len(drivers))
}
