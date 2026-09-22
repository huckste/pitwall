package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"pitwall/backend/db"
	"pitwall/backend/models"
	"pitwall/backend/ocb"
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

	teams, err := ocb.FetchSeasonTeams("formula1", season2025.ID)

	if err != nil {
		log.Fatal(err)
	}

	for _, t := range teams {
		if err := db.UpsertTeam(ctx, pool, models.MapSeasonTeam(t, models.Formula_1)); err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("upserted %d teams for formula1 2025\n", len(teams))
}
