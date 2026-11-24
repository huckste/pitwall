package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DriverCard is a driver's season at a glance, for the player card view.
// The season is the one their most recent race belongs to.
type DriverCard struct {
	DriverID    int     `json:"driverId"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Series      string  `json:"series"`
	Number      *int    `json:"number"`
	Nationality *string `json:"nationality"`
	Team        *string `json:"team"`
	TeamColor   *string `json:"teamColor"`
	// PhotoURL is a Wikimedia Commons image; PhotoPage is its Wikipedia page, for attribution.
	PhotoURL   *string      `json:"photoUrl"`
	PhotoPage  *string      `json:"photoPage"`
	SeasonYear int          `json:"seasonYear"`
	ChampPos   int          `json:"champPos"`
	ChampOf    int          `json:"champOf"`
	Points     float64      `json:"points"`
	Starts     int          `json:"starts"`
	Wins       int          `json:"wins"`
	Podiums    int          `json:"podiums"`
	BestFinish *int         `json:"bestFinish"`
	AvgFinish  *float64     `json:"avgFinish"`
	DNFs       int          `json:"dnfs"`
	RecentForm []FormResult `json:"recentForm"`
	Career     Career       `json:"career"`
	// Rarity is the card tier: champion, winner, rookie or base. See cardRarity.
	Rarity string `json:"rarity"`
}

// Career is everything we have stored for the driver in this series, across seasons.
// It only goes back as far as the seasons that have been backfilled.
type Career struct {
	Seasons    int   `json:"seasons"`
	Starts     int   `json:"starts"`
	Wins       int   `json:"wins"`
	Podiums    int   `json:"podiums"`
	Titles     int   `json:"titles"`
	TitleYears []int `json:"titleYears"`
	IsRookie   bool  `json:"isRookie"`
}

// cardRarity is the one place the card tier rules live.
// Priority order matters: a rookie who has won a race is a "winner" card
// (the rookie badge still shows separately via Career.IsRookie).
func cardRarity(c Career) string {
	switch {
	case c.Titles > 0:
		return "champion"
	case c.Wins > 0:
		return "winner"
	case c.IsRookie:
		return "rookie"
	default:
		return "base"
	}
}

// FormResult is one race in the driver's recent form strip.
type FormResult struct {
	Event    string  `json:"event"`
	Position *string `json:"position"`
	Outcome  string  `json:"outcome"`
}

// GetDriverCard returns nil (and no error) when the driver has no race results.
func GetDriverCard(ctx context.Context, pool *pgxpool.Pool, driverID int) (*DriverCard, error) {
	var c DriverCard
	var seasonID int

	err := pool.QueryRow(ctx, `
		with target as (
			-- the driver's most recent race decides which season the card shows
			select d.id as driver_id, e.season_id, r.team_id
			from drivers d
			join results r  on r.driver_id = d.id
			join sessions s on s.id = r.session_id
			join events e   on e.id = s.event_id
			where d.id = $1 and s.type = 'race'
			order by s.start_time desc
			limit 1
		),
		season_results as (
			-- every result in that season, for every driver (needed for championship rank)
			select
				r.driver_id,
				r.points,
				s.type,
				`+outcomeSQL+` as outcome,
				case when r.position ~ '^\d+$' then r.position::int end as pos_num
			from results r
			join sessions s on s.id = r.session_id
			join events e   on e.id = s.event_id
			where e.season_id = (select season_id from target)
		),
		totals as (
			select
				driver_id,
				coalesce(sum(points), 0)::float8                                              as points,
				count(*) filter (where type = 'race')                                         as starts,
				count(*) filter (where type = 'race' and outcome = 'finished' and pos_num = 1)  as wins,
				count(*) filter (where type = 'race' and outcome = 'finished' and pos_num <= 3) as podiums,
				min(pos_num) filter (where type = 'race' and outcome = 'finished')            as best_finish,
				avg(pos_num) filter (where type = 'race' and outcome = 'finished')::float8    as avg_finish,
				count(*) filter (where type = 'race' and outcome <> 'finished')               as dnfs
			from season_results
			group by driver_id
		),
		ranked as (
			-- only drivers who started a race count toward the championship table
			select *,
				rank() over (order by points desc) as champ_pos,
				count(*) over ()                   as champ_of
			from totals
			where starts > 0
		)
		select
			d.id, d.first_name, d.last_name, d.series, d.number, d.nationality,
			tm.name, tm.color, d.photo_url, d.photo_page, se.year, t.season_id,
			k.champ_pos, k.champ_of, k.points, k.starts, k.wins, k.podiums,
			k.best_finish, k.avg_finish, k.dnfs
		from target t
		join drivers d    on d.id = t.driver_id
		join ranked k     on k.driver_id = t.driver_id
		join seasons se   on se.id = t.season_id
		left join teams tm on tm.id = t.team_id
	`, driverID).Scan(
		&c.DriverID, &c.FirstName, &c.LastName, &c.Series, &c.Number, &c.Nationality,
		&c.Team, &c.TeamColor, &c.PhotoURL, &c.PhotoPage, &c.SeasonYear, &seasonID,
		&c.ChampPos, &c.ChampOf, &c.Points, &c.Starts, &c.Wins, &c.Podiums,
		&c.BestFinish, &c.AvgFinish, &c.DNFs,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	form, err := recentForm(ctx, pool, driverID, seasonID)
	if err != nil {
		return nil, err
	}
	c.RecentForm = form

	career, err := careerStats(ctx, pool, driverID, c.Series, c.SeasonYear)
	if err != nil {
		return nil, err
	}
	c.Career = career
	c.Rarity = cardRarity(career)

	return &c, nil
}

// careerStats totals the driver's races across every stored season of the series,
// counts championships, and works out whether this is their rookie season.
func careerStats(ctx context.Context, pool *pgxpool.Pool, driverID int, series string, seasonYear int) (Career, error) {
	var c Career
	var driverFirstYear, seriesFirstYear *int

	err := pool.QueryRow(ctx, `
		with driver_races as (
			select
				e.season_id,
				se.year,
				`+outcomeSQL+` as outcome,
				case when r.position ~ '^\d+$' then r.position::int end as pos_num
			from results r
			join sessions s on s.id = r.session_id
			join events e   on e.id = s.event_id
			join seasons se on se.id = e.season_id
			where r.driver_id = $1 and s.type = 'race'
		),
		series_seasons as (
			-- a season only counts as decided once none of its sessions are still scheduled;
			-- without this the current leader would be shown as champion mid-season
			select e.season_id, se.year, bool_and(s.status <> 'scheduled') as complete
			from events e
			join seasons se on se.id = e.season_id
			join sessions s on s.event_id = e.id
			where e.series = $2
			group by e.season_id, se.year
		),
		season_points as (
			select e.season_id, r.driver_id, sum(r.points) as pts
			from results r
			join sessions s on s.id = r.session_id
			join events e   on e.id = s.event_id
			where e.series = $2
			group by e.season_id, r.driver_id
		),
		champions as (
			select distinct on (sp.season_id) sp.season_id, sp.driver_id
			from season_points sp
			join series_seasons ss on ss.season_id = sp.season_id and ss.complete
			order by sp.season_id, sp.pts desc nulls last
		)
		select
			count(distinct season_id),
			count(*),
			count(*) filter (where outcome = 'finished' and pos_num = 1),
			count(*) filter (where outcome = 'finished' and pos_num <= 3),
			(select coalesce(array_agg(ss.year order by ss.year), '{}')
			 from champions ch join series_seasons ss on ss.season_id = ch.season_id
			 where ch.driver_id = $1),
			min(year),
			(select min(year) from series_seasons)
		from driver_races
	`, driverID, series).Scan(
		&c.Seasons, &c.Starts, &c.Wins, &c.Podiums, &c.TitleYears,
		&driverFirstYear, &seriesFirstYear,
	)
	if err != nil {
		return Career{}, err
	}
	c.Titles = len(c.TitleYears)

	// Rookie = first season we have for them is this one, and we have an earlier
	// season of the series to compare against. Without the second check, every
	// driver in a series with only one stored season would show as a rookie.
	if driverFirstYear != nil && seriesFirstYear != nil {
		c.IsRookie = *driverFirstYear == seasonYear && *driverFirstYear > *seriesFirstYear
	}

	return c, nil
}

// recentForm returns the driver's last 5 races in the season, newest first.
func recentForm(ctx context.Context, pool *pgxpool.Pool, driverID, seasonID int) ([]FormResult, error) {
	rows, err := pool.Query(ctx, `
		select e.name, r.position, `+outcomeSQL+` as outcome
		from results r
		join sessions s on s.id = r.session_id
		join events e   on e.id = s.event_id
		where r.driver_id = $1 and e.season_id = $2 and s.type = 'race'
		order by s.start_time desc
		limit 5
	`, driverID, seasonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	form := []FormResult{}
	for rows.Next() {
		var f FormResult
		if err := rows.Scan(&f.Event, &f.Position, &f.Outcome); err != nil {
			return nil, err
		}
		form = append(form, f)
	}
	return form, rows.Err()
}
