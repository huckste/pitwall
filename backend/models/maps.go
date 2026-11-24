package models

import (
	"pitwall/backend/db"
	"pitwall/backend/ocb"
	"time"
)

const (
	Formula1      Series = "formula1"
	NASCAR        Series = "nascar"
	NASCARXfinity Series = "nascar-xfinity"
	NASCARTrucks  Series = "nascar-truck"
	IndyCar       Series = "indycar"
	MotoGP        Series = "moto-gp"
	Moto2         Series = "moto2"
	Moto3         Series = "moto3"
	FormulaE      Series = "formula-e"
	Formula2      Series = "formula2"
	Formula3      Series = "formula3"
	WEC           Series = "wec"
	WRC           Series = "wrc"
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

func MapSeasonDriver(d ocb.SeasonDriver, series Series) db.Driver {
	var teamID *string
	best := 0

	// Driver can be an many teams within a season so this grabs the team ID that contains hightest round value
	for _, t := range d.Teams {
		for _, round := range t.ParticipationRounds {
			if round > best {
				best = round
				teamID = &t.ID
			}
		}
	}

	return db.Driver{
		TeamOCBID:   teamID,
		OCBID:       d.ID,
		Series:      string(series),
		FirstName:   d.FirstName,
		LastName:    d.LastName,
		Number:      d.Number,
		Nationality: d.Country.Name,
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

func MapSeason(s ocb.Season, series Series) db.Season {
	return db.Season{
		OCBID:      s.ID,
		Series:     string(series),
		Year:       s.Year,
		RoundCount: s.RoundCount,
		Status:     s.Status,
	}
}

func MapEvent(e ocb.Event, series Series, seasonOCBID string) (db.Event, error) {
	start, err := time.Parse("2006-01-02", e.DateStart)
	if err != nil {
		return db.Event{}, err
	}

	end, err := time.Parse("2006-01-02", e.DateEnd)
	if err != nil {
		return db.Event{}, err
	}

	return db.Event{
		OCBID:       e.ID,
		Series:      string(series),
		Name:        e.Name,
		Status:      e.Status,
		DateStart:   start,
		DateEnd:     end,
		SeasonOCBID: seasonOCBID,
	}, nil
}

func MapSession(s ocb.ScheduleSession, series Series, eventOCBID string) (db.Session, error) {
	var start *time.Time
	var end *time.Time

	if s.StartTime != "" {
		t, err := time.Parse(time.RFC3339, s.StartTime)
		if err != nil {
			return db.Session{}, err
		}
		start = &t
	}

	if s.EndTime != "" {
		t, err := time.Parse(time.RFC3339, s.EndTime)
		if err != nil {
			return db.Session{}, err
		}
		end = &t
	}

	return db.Session{
		OCBID:      s.ID,
		Type:       s.Type,
		Name:       s.Name,
		Status:     s.Status,
		Series:     string(series),
		EventOCBID: eventOCBID,
		StartTime:  start,
		EndTime:    end,
	}, nil
}
