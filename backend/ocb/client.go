// Package ocb is a client for the Orange Cat Blacktop racing data API
package ocb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const baseURL = "https://api.ocblacktop.com/v1"

func get(url string, out any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("x-api-key", os.Getenv("OCB_API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// ignore error given during closing as there is nothing to do with the err
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ocb api returned %d: %s/%s", resp.StatusCode, url, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return err
	}

	return nil
}

func FetchDrivers(sport string) ([]Driver, error) {
	url := fmt.Sprintf("%s/%s/drivers", baseURL, sport)
	var parsed driversResponse

	if err := get(url, &parsed); err != nil {
		return nil, err
	}

	return parsed.Data, nil
}

func FetchEvents(sport string) ([]Event, error) {
	url := fmt.Sprintf("%s/%s/events", baseURL, sport)
	var parsed eventsResponse

	if err := get(url, &parsed); err != nil {
		return nil, err
	}

	return parsed.Data, nil
}

func MostRecentCompletedEvent(events []Event) (*Event, bool) {
	for _, e := range events {
		if e.Status == "completed" {
			return &e, true
		}
	}

	return nil, false
}

func FetchSeasons(sport string) ([]Season, error) {
	url := fmt.Sprintf("%s/%s/seasons?page=1&limit=10", baseURL, sport)
	var parsed seasonsResponse

	if err := get(url, &parsed); err != nil {
		return nil, err
	}

	return parsed.Data, nil
}

func FetchSeasonTeams(sport, seasonID string) ([]SeasonTeam, error) {
	url := fmt.Sprintf("%s/%s/seasons/%s/teams", baseURL, sport, seasonID)
	var teams []SeasonTeam

	if err := get(url, &teams); err != nil {
		return nil, err
	}

	return teams, nil
}
