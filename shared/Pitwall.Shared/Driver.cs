namespace Pitwall.Shared;

public class Driver
{
  public required string Id { get; set; }
  public required string FirstName { get; set; }
  public required string LastName { get; set; }
  public required Series Series { get; set; }
  public int? DriverNumber { get; set; }
  public string? TeamName { get; set; }
  public string? Nationality { get; set; }
}
