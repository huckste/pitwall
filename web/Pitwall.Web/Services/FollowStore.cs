using System.Text.Json;
using Microsoft.JSInterop;

namespace Pitwall.Web.Services;

// Keeps the set of followed drivers and the chosen view mode in the browser's
// localStorage. Temporary home until user accounts exist; then this moves to the API.
public class FollowStore(IJSRuntime js)
{
    private const string FollowsKey = "pitwall.follows";
    private const string ViewKey = "pitwall.view";

    public HashSet<int> DriverIds { get; private set; } = new();
    public string View { get; private set; } = "consolidated";

    public async Task LoadAsync()
    {
        var follows = await GetAsync(FollowsKey);
        if (follows is not null)
        {
            try
            {
                DriverIds = JsonSerializer.Deserialize<HashSet<int>>(follows) ?? new();
            }
            catch (JsonException)
            {
                DriverIds = new();
            }
        }

        View = await GetAsync(ViewKey) ?? "consolidated";
    }

    public bool IsFollowed(int driverId) => DriverIds.Contains(driverId);

    public async Task ToggleAsync(int driverId)
    {
        if (!DriverIds.Remove(driverId))
        {
            DriverIds.Add(driverId);
        }
        await SetAsync(FollowsKey, JsonSerializer.Serialize(DriverIds));
    }

    public async Task SetViewAsync(string view)
    {
        View = view;
        await SetAsync(ViewKey, view);
    }

    // localStorage can be blocked (private mode, disabled storage); the app
    // still works for the session, it just won't remember anything.
    private async Task<string?> GetAsync(string key)
    {
        try
        {
            return await js.InvokeAsync<string?>("localStorage.getItem", key);
        }
        catch (JSException)
        {
            return null;
        }
    }

    private async Task SetAsync(string key, string value)
    {
        try
        {
            await js.InvokeVoidAsync("localStorage.setItem", key, value);
        }
        catch (JSException)
        {
        }
    }
}
