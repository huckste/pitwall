namespace Pitwall.Shared;

public class RaceResults
{
  public required string RaceId { get; set; }
  public required string DriverId { get; set; }
  public int? FinishPosition { get; set; }
  public int? StartPosition { get; set; }
  public double Points { get; set; }
  public required string Status { get; set; }
}
