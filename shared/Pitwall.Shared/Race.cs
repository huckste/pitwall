namespace Pitwall.Shared;

public class Race
{
  public required string Id { get; set; }
  public required Series Series { get; set; }
  public required string Name { get; set; }
  public required DateTime Date { get; set; }
  public int Round { get; set; }
}
