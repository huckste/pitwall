package models

import "time"

type Series string

type Driver struct {
	ID           string  `json:"id"`
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	Series       Series  `json:"series"`
	DriverNumber *int    `json:"driverNumber,omitempty"`
	TeamName     *string `json:"teamName,omitempty"`
	Nationality  *string `json:"nationality,omitempty"`
}

type Race struct {
	ID     string    `json:"id"`
	Series Series    `json:"series"`
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`
	Round  int       `json:"round"`
}

type RaceResults struct {
	RaceID         string  `json:"raceId"`
	DriverID       string  `json:"driverId"`
	FinishPosition *int    `json:"finishPosition,omitempty"`
	StartPosition  *int    `json:"startPosition,omitempty"`
	Points         float64 `json:"points"`
	Status         string  `json:"status"`
}
