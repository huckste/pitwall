package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"pitwall/backend/models"
	"pitwall/backend/ocb"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/drivers", driversHandler)
	http.HandleFunc("/latest-race", latestRaceHandler)

	log.Println("pitwall backend listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func driversHandler(w http.ResponseWriter, r *http.Request) {

	seriesParam := models.Series(r.URL.Query().Get("series"))
	ocbDrivers, err := ocb.FetchDrivers(string(seriesParam))

	if err != nil {
		http.Error(w, "failed to fetch drivers", http.StatusBadGateway)
		return
	}

	drivers := make([]models.Driver, 0, len(ocbDrivers))

	for _, d := range ocbDrivers {
		drivers = append(drivers, models.MapOCBDriver(d, seriesParam))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drivers)
}

func latestRaceHandler(w http.ResponseWriter, r *http.Request) {

	seriesParam := models.Series(r.URL.Query().Get("series"))
	events, err := ocb.FetchEvents(string(seriesParam))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	event, found := ocb.MostRecentCompletedEvent(events)

	if !found {
		http.Error(w, "no completed race found", http.StatusNotFound)
		return
	}

	race, err := models.MapOCBEvent(*event, seriesParam)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(race)
}
