namespace Pitwall.Shared;

public record LatestResult(
    int DriverId,
    string FirstName,
    string LastName,
    string Series,
    string? Team,
    string Event,
    DateTimeOffset? StartTime,
    string? Position,
    double? Points,
    string Status,
    string Outcome
    );
