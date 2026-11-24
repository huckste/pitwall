namespace Pitwall.Shared;

// PhotoUrl is a Wikimedia Commons image found via Wikipedia; PhotoPage is that page (for attribution).
// A driver's card (GET /drivers/{id}/card): current season + career in this series.
public record DriverCard(
    int DriverId,
    string FirstName,
    string LastName,
    string Series,
    int? Number,
    string? Nationality,
    string? Team,
    string? TeamColor,
    string? PhotoUrl,
    string? PhotoPage,
    int SeasonYear,
    int ChampPos,
    int ChampOf,
    double Points,
    int Starts,
    int Wins,
    int Podiums,
    int? BestFinish,
    double? AvgFinish,
    int DNFs,
    List<FormResult> RecentForm,
    Career Career,
    string Rarity);

// Totals across every stored season of the series (only as far back as the backfill goes).
public record Career(int Seasons, int Starts, int Wins, int Podiums, int Titles, List<int> TitleYears, bool IsRookie);

// One race in the recent-form strip, newest first.
public record FormResult(string Event, string? Position, string Outcome);
