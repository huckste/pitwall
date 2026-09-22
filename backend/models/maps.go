package models

import (
	"pitwall/backend/db"
	"pitwall/backend/ocb"
	"time"
)

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

func MapOCBDriver(d ocb.Driver, series Series) Driver {

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

func MapOCBEvent(e ocb.Event, series Series) (Race, error) {

	date, err := time.Parse("2006-01-02", e.DateStart)

	if err != nil {
		return Race{}, err
	}

	return Race{
		ID:     e.ID,
		Series: series,
		Name:   e.Name,
		Date:   date,
		Round:  0,
	}, nil

}

func MapSeasonTeam(t ocb.SeasonTeam, series Series) db.Team {
	shortName := t.Name
	if t.ShortName != nil {
		shortName = *t.ShortName
	}

	fullName := t.Name
	if t.FullName != nil {
		fullName = *t.FullName
	}

	return db.Team{
		OCBID:     t.ID,
		Series:    string(series),
		Name:      t.Name,
		ShortName: shortName,
		FullName:  fullName,
	}
}
