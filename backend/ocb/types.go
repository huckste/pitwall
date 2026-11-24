package ocb

type Location struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	City    *string `json:"city"`
	Country Country `json:"country"`
}

type Season struct {
	ID         string `json:"id"`
	Year       int    `json:"year"`
	SportID    string `json:"sportId"`
	Status     string `json:"status"`
	RoundCount *int   `json:"roundCount"`
}

type ScheduleSession struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Status    string `json:"status"`
}

type Event struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	DateStart string            `json:"dateStart"`
	DateEnd   string            `json:"dateEnd"`
	Status    string            `json:"status"`
	Location  Location          `json:"location"`
	SportID   string            `json:"sportId"`
	Schedule  []ScheduleSession `json:"schedule"`
}

type Country struct {
	Name      *string `json:"name"`
	TwoCode   *string `json:"twoCode"`
	ThreeCode *string `json:"threeCode"`
}

type Driver struct {
	ID        string  `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	BirthDate string  `json:"birthDate"`
	Number    *int    `json:"number"`
	TLA       *string `json:"tla"`
	Country   Country `json:"country"`
}

type SeasonDriver struct {
	ID        string             `json:"id"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Number    *int               `json:"number"`
	Country   Country            `json:"country"`
	Teams     []SeasonDriverTeam `json:"teams"`
}

type SeasonDriverTeam struct {
	ID                  string `json:"id"`
	ParticipationRounds []int  `json:"participationRounds"`
}

type SeasonTeam struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	ShortName *string `json:"shortName"`
	FullName  *string `json:"fullName"`
	Color     *string `json:"color"`
}

type driversResponse struct {
	Data []Driver `json:"data"`
}

type eventsResponse struct {
	Data []Event `json:"data"`
}

type seasonsResponse struct {
	Data []Season `json:"data"`
}

type ResultDriver struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Number    *int   `json:"number"`
}

type ResultRef struct {
	ID string `json:"id"`
}

type SessionResult struct {
	Position  *string      `json:"position"`
	Points    *string      `json:"points"`
	Laps      *int         `json:"laps"`
	Status    string       `json:"status"`
	CarNumber *string      `json:"carNumber"`
	Driver    ResultDriver `json:"driver"`
	Team      ResultRef    `json:"team"`
}

type SessionDetail struct {
	Results []SessionResult `json:"results"`
}
