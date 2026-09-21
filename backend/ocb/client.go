package ocb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const ocbBaseURL = "https://api.ocblacktop.com/v1"

type ocbCountry struct {
	Name      *string `json:"name"`
	TwoCode   *string `json:"twoCode"`
	ThreeCode *string `json:"threeCode"`
}

type ocbDriver struct {
	ID        string     `json:"id"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	BirthDate string     `json:"birthDate"`
	Number    *int       `json:"number"`
	TLA       *string    `json:"tla"`
	Country   ocbCountry `json:"country"`
}

type ocbDriversResponse struct {
	Data []ocbDriver `json:"data"`
}

func FetchDrivers(sport string) ([]ocbDriver, error) {
	url := fmt.Sprintf("%s/%s/drivers", ocbBaseURL, sport)

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

	var parsed ocbDriversResponse

	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	return parsed.Data, nil
}
