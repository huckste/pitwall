package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// outcomeSQL turns each series' own status words into one shared vocabulary:
// finished, dnf, dns, dsq. Status words differ per series (F1 "OK",
// F2/F3/FE "", IndyCar "Running"), so this is the one place that knows them.
// Expects the results table to be aliased as r.
const outcomeSQL = `case
	when r.status in ('', 'OK', 'Running') then 'finished'
	when r.status = 'DNS' then 'dns'
	when r.status = 'DSQ' then 'dsq'
	else 'dnf'
end`

type LatestResult struct {
	DriverID  int        `json:"driverId"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Series    string     `json:"series"`
	Team      *string    `json:"team"`
	Event     string     `json:"event"`
	StartTime *time.Time `json:"startTime"`
	Position  *string    `json:"position"`
	Points    *float64   `json:"points"`
	Status    string     `json:"status"`
	Outcome   string     `json:"outcome"`
}

func LatestRaceResult(ctx context.Context, pool *pgxpool.Pool) ([]LatestResult, error) {
	rows, err := pool.Query(ctx, `
		select distinct on (d.id)
			d.id,
			d.first_name,
			d.last_name,
			d.series,
			t.name as team, 
			e.name as event,
			s.start_time,
			r.position,
			r.points,
			r.status,
			`+outcomeSQL+` as outcome
		from drivers d
		join results r 
			on r.driver_id = d.id
		join sessions s 
			on s.id = r.session_id
		join events e 
			on e.id = s.event_id
		left join teams t 
			on t.id = r.team_id
		where s.type = 'race'
		order by 
			d.id, 
			s.start_time desc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []LatestResult
	for rows.Next() {
		var lr LatestResult
		if err := rows.Scan(&lr.DriverID, &lr.FirstName, &lr.LastName, &lr.Series,
			&lr.Team, &lr.Event, &lr.StartTime, &lr.Position, &lr.Points, &lr.Status, &lr.Outcome); err != nil {
			return nil, err
		}
		out = append(out, lr)
	}

	return out, rows.Err()
}

func SessionResultsFetched(ctx context.Context, pool *pgxpool.Pool, sessionOCBID string) (bool, error) {
	var fetched bool
	err := pool.QueryRow(ctx, `
		SELECT results_fetched_at IS NOT NULL FROM sessions WHERE ocb_id = $1
	`, sessionOCBID).Scan(&fetched)
	return fetched, err
}
