package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"pitwall/backend/db"
	"pitwall/backend/ingest"
	"pitwall/backend/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatal("failed to connect to db: ", err)
	}
	defer pool.Close()

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/latest-results", latestResultHandler(pool))

	log.Println("pitwall backend listening on :8080")

	tracked := []models.Series{models.Formula1, models.Formula2, models.Formula3, models.IndyCar, models.FormulaE}

	go func() {
		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()

		for {
			for _, s := range tracked {
				log.Printf("poller: syncing %s", s)
				if err := ingest.SyncSeason(ctx, pool, s, time.Now().Year()); err != nil {
					log.Printf("poller: %s sync failed: %v", s, err)
				}
			}

			<-ticker.C
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", withCORS(http.DefaultServeMux)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func latestResultHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results, err := db.LatestRaceResult(r.Context(), pool)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
