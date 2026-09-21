package main

import (
	"pitwall/backend/ocb"
	"time"
)

type Series string

const (
	Formula_1      Series = "formula1"
	NASCAR         Series = "nascar"
	NASCAR_Xfinity Series = "nascar-xfinity"
	NASCAR_Trucks  Series = "nascar-truck"
	IndyCar        Series = "indycar"
	Moto_GP        Series = "moto-gp"
	Moto_2         Series = "moto2"
	Moto_3         Series = "moto3"
	Formula_E      Series = "formula-e"
	Formula_2      Series = "formula2"
	Formula_3      Series = "formula3"
	WEC            Series = "wec"
	WRC            Series = "wrc"
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

func mapOCBDriver(d ocb.Driver, series Series) Driver {

	var nationality *string

	if d.Country.Name != nil {
		nationality = d.Country.Name
	}

	return Driver{
		ID:           d.ID,
		FirstName:    d.FirstName,
		LastName:     d.LastName,
		Series:       series,
		DriverNumber: d.Number,
		TeamName:     nil,
		Nationality:  nationality,
	}
}
