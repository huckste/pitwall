namespace Pitwall.Web.Services;

// Display details for each OCB series slug. Keeps slug knowledge in one place.
public static class SeriesInfo
{
    public static string Name(string slug) => slug switch
    {
        "formula1" => "Formula 1",
        "formula2" => "Formula 2",
        "formula3" => "Formula 3",
        "formula-e" => "Formula E",
        "indycar" => "IndyCar",
        "wec" => "WEC",
        "nascar" => "NASCAR Cup",
        _ => slug
    };

    public static string Short(string slug) => slug switch
    {
        "formula1" => "F1",
        "formula2" => "F2",
        "formula3" => "F3",
        "formula-e" => "FE",
        "indycar" => "INDY",
        "wec" => "WEC",
        "nascar" => "CUP",
        _ => slug.ToUpper()
    };

    public static string Accent(string slug) => slug switch
    {
        "formula1" => "var(--accent-1)",
        "formula2" => "var(--accent-2)",
        "formula3" => "var(--accent-3)",
        "formula-e" => "var(--accent-4)",
        "indycar" => "var(--accent-5)",
        _ => "var(--text-muted)"
    };
}
