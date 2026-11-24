package main

import (
	"context"
	"flag"
	"log"
	"pitwall/backend/db"
	"pitwall/backend/ingest"
	"pitwall/backend/models"

	"github.com/joho/godotenv"
)

func main() {
	seriesFlag := flag.String("series", "formula1", "OCB series slug")
	yearFlag := flag.Int("year", 2025, "season year")
	flag.Parse()

	series := models.Series(*seriesFlag)
	year := *yearFlag

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

	if err := ingest.SyncSeason(ctx, pool, series, year); err != nil {
		log.Fatal(err)
	}
}
