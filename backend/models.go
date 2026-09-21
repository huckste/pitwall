package main

import "time"

type Series string

const (
	Formula_1      Series = "Formula1"
	NASCAR         Series = "Nascar"
	NASCAR_Xfinity Series = "NascarXfinity"
	NASCAR_Trucks  Series = "NascarTrucks"
	IndyCar        Series = "Indycar"
	Moto_GP        Series = "MotoGp"
	Moto_2         Series = "Moto2"
	Moto_3         Series = "Moto3"
	Formula_E      Series = "FormulaE"
	Formula_2      Series = "Formula2"
	Formula_3      Series = "Formula3"
	WEC            Series = "Wec"
	WRC            Series = "Wrc"
)

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
