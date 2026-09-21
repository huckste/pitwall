package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type openF1Driver struct {
	DriverNumber int    `json:"driver_number"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	FullName     string `json:"full_name"`
	TeamName     string `json:"team_name"`
	CountryCode  string `json:"countery_code"`
}

func main() {
	resp, err := http.Get("https://api.openf1.org/v1/drivers?session_key=latest")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	var drivers []openF1Driver

	if err := json.NewDecoder(resp.Body).Decode(&drivers); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Parsed %d drivers\n", len(drivers))

	for _, d := range drivers {
		fmt.Printf("#%d %s (%s)\n", d.DriverNumber, d.FullName, d.TeamName)
	}

}
