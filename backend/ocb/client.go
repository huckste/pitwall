package ocb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const baseURL = "https://api.ocblacktop.com/v1"

func FetchDrivers(sport string) ([]Driver, error) {
	url := fmt.Sprintf("%s/%s/drivers", baseURL, sport)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", os.Getenv("OCB_API_KEY"))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var parsed driversResponse

	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	return parsed.Data, nil
}

func FetchEvents(sport string) ([]Event, error) {
	url := fmt.Sprintf("%s/%s/events", baseURL, sport)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", os.Getenv("OCB_API_KEY"))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ocb api returned %d: %s", resp.StatusCode, string(body))
	}

	var parsed eventsResponse

	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
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

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", os.Getenv("OCB_API_KEY"))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ocb api returned %d: %s", resp.StatusCode, string(body))
	}

	var parsed seasonsResponse

	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	return parsed.Data, nil
}

func FetchSeasonTeams(sport, seasonID string) ([]SeasonTeam, error) {

	url := fmt.Sprintf("%s/%s/seasons/%s/teams", baseURL, sport, seasonID)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", os.Getenv("OCB_API_KEY"))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ocb api returned %d: %s", resp.StatusCode, string(body))
	}

	var teams []SeasonTeam

	if err := json.NewDecoder(resp.Body).Decode(&teams); err != nil {
		return nil, err
	}

	return teams, nil
}
