package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Team struct {
	OCBID     string
	Series    string
	Name      string
	ShortName string
	FullName  string
}

type Driver struct {
	OCBID       string
	Series      string
	FirstName   string
	LastName    string
	Number      *int
	Nationality *string
	TeamOCBID   *string
}

type Season struct {
	OCBID      string
	Series     string
	Year       int
	RoundCount *int
	Status     string
}

type Event struct {
	OCBID       string
	Series      string
	Name        string
	Status      string
	DateStart   time.Time
	DateEnd     time.Time
	SeasonOCBID string
}

type Session struct {
	OCBID      string
	Type       string
	Name       string
	Status     string
	StartTime  *time.Time
	EndTime    *time.Time
	EventOCBID string
	Series     string
}

type Result struct {
	SessionOCBID string
	DriverOCBID  string
	TeamOCBID    string
	Series       string
	Position     *string
	CarNumber    *string
	Points       *float64
	Laps         *int
	Status       string
}

func UpsertTeam(ctx context.Context, pool *pgxpool.Pool, t Team) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO teams (ocb_id, series, name, short_name, full_name)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (ocb_id, series) DO UPDATE
		SET name = EXCLUDED.name, 
			short_name = EXCLUDED.short_name, 
			full_name = EXCLUDED.full_name
	`, t.OCBID, t.Series, t.Name, t.ShortName, t.FullName)

	return err
}

func UpsertDriver(ctx context.Context, pool *pgxpool.Pool, d Driver) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO drivers (ocb_id, series, first_name, last_name, number, nationality, team_id)
		VALUES ($1, $2, $3, $4, $5, $6,
			(SELECT id FROM teams WHERE ocb_id = $7 AND series = $2))
		ON CONFLICT (ocb_id, series) DO UPDATE
		SET first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			number = EXCLUDED.number,
			nationality = EXCLUDED.nationality,
			team_id = EXCLUDED.team_id
		`, d.OCBID, d.Series, d.FirstName, d.LastName, d.Number, d.Nationality, d.TeamOCBID)

	return err
}

func UpsertSeason(ctx context.Context, pool *pgxpool.Pool, s Season) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO seasons (ocb_id, series, year, round_count, status) 
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (ocb_id, series) DO UPDATE
		SET year = EXCLUDED.year,
			round_count = EXCLUDED.round_count,
			status = EXCLUDED.status
		`, s.OCBID, s.Series, s.Year, s.RoundCount, s.Status)

	return err
}

func UpsertEvent(ctx context.Context, pool *pgxpool.Pool, e Event) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO events (ocb_id, series, name, date_start, date_end, status, season_id )
		VALUES ($1, $2, $3, $4, $5, $6,
			(SELECT id FROM seasons WHERE ocb_id = $7 AND series = $2))
		ON CONFLICT (ocb_id, series) DO UPDATE
		SET name = EXCLUDED.name,
			date_start = EXCLUDED.date_start,
			date_end = EXCLUDED.date_end,
			status = EXCLUDED.status,
			season_id = EXCLUDED.season_id
		`, e.OCBID, e.Series, e.Name, e.DateStart, e.DateEnd, e.Status, e.SeasonOCBID)

	return err
}

func UpsertSession(ctx context.Context, pool *pgxpool.Pool, s Session) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO sessions (ocb_id, type, name, status, start_time, end_time, event_id)
		VALUES ($1, $2, $3, $4, $5, $6,
			(SELECT id FROM events WHERE ocb_id = $7 AND series = $8)) 
		ON CONFLICT (ocb_id) DO UPDATE
		SET name = EXCLUDED.name,
			start_time = EXCLUDED.start_time,
			end_time = EXCLUDED.end_time,
			event_id = EXCLUDED.event_id,
			type = EXCLUDED.type,
			status = EXCLUDED.status
		`, s.OCBID, s.Type, s.Name, s.Status, s.StartTime, s.EndTime, s.EventOCBID, s.Series)

	return err
}

func UpsertResult(ctx context.Context, pool *pgxpool.Pool, r Result) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO results (session_id, driver_id, team_id, car_number, position, points, status, laps)
		VALUES (
			(SELECT id FROM sessions WHERE ocb_id = $1),
			(SELECT id FROM drivers WHERE ocb_id = $2 AND series = $3),
			(SELECT id FROM teams WHERE ocb_id = $4 AND series = $3),
			$5, $6, $7, $8, $9
		)
		ON CONFLICT (session_id, driver_id) DO UPDATE 
		SET team_id = EXCLUDED.team_id,
			car_number = EXCLUDED.car_number,
			position = EXCLUDED.position,
			points = EXCLUDED.points,
			status = EXCLUDED.status,
			laps = EXCLUDED.laps
		`, r.SessionOCBID, r.DriverOCBID, r.Series, r.TeamOCBID, r.CarNumber, r.Position, r.Points, r.Status, r.Laps)

	return err
}

func InsertMissingDriver(ctx context.Context, pool *pgxpool.Pool, d Driver) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO drivers (ocb_id, series, first_name, last_name, number, nationality, team_id)
		VALUES ($1, $2, $3, $4, $5, $6, 
			(SELECT id FROM teams WHERE ocb_id = $7 AND series = $2))
		ON CONFLICT (ocb_id, series) DO NOTHING
		`, d.OCBID, d.Series, d.FirstName, d.LastName, d.Number, d.Nationality, d.TeamOCBID)

	return err
}
