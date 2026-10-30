# OCB API notes

Findings from live testing of the Orange Cat Blacktop API. The docs and real responses have disagreed more than once, so where this file and the spec differ, trust this file.

## Quick reference

| Item             | Value                                                                                    |
| ---------------- | ---------------------------------------------------------------------------------------- |
| Base URL         | `https://api.ocblacktop.com/v1`                                                          |
| Auth             | `x-api-key` header                                                                       |
| Free tier limits | 7,500 requests/month, 60/minute                                                          |
| Full spec        | `~/Downloads/api-1.json` (not in this repo)                                              |
| Series covered   | F1, F2, F3, Formula E, IndyCar, NASCAR (Cup/Truck/Xfinity), MotoGP/Moto2/Moto3, WEC, WRC |
| Not covered      | **IMSA**: not available anywhere in this API                                             |

## Top gotchas

Read these first. Each one has more detail further down.

1. **URL slug ≠ `sportId`.** You call `formula1`, but responses say `"formula-1"`.
2. **Same field, different types.** `points` and `position` are strings in some endpoints and numbers in others.
3. **Status values differ per sport.** There's no shared enum.
4. **The docs' example status values are wrong.** A real season's status is `"Created"`.
5. **WEC standings are per crew, not per driver.** Sum session results yourself.
6. **NASCAR stages have no driver ID.** Only name matching works.

## URL slug vs `sportId`

- **In the URL path:** a lowercase slug, dashed only for compound names: `formula1`, `nascar`, `nascar-trucks`, `indycar`, `wec`.
- **In response bodies (`sportId`):** always the dashed form, e.g. `"formula-1"`, even when the URL slug has no dash.
- **Rule:** never compare one against the other or reuse one as the other.

## Endpoints

All paths follow the pattern `/{sport}/...`.

### Drivers: `/{sport}/drivers`, `/{sport}/drivers/{id}`

**Use for:** driver lists and profiles.

**Confirmed identical shape across:** F1, Formula E, F2.

**Fields:** `id`, `firstName`, `lastName`, `birthDate`, `number`, `tla`, `country` (`name`, `twoCode`, `threeCode`)

### Teams: `/{sport}/teams`, `/{sport}/teams/{id}`

**Use for:** following a team or manufacturer.

- Team detail includes a `drivers[]` roster.
- Each roster driver has `teamSeasonParticipations[]`.
- Team detail also has `seasonParticipations[]`.

### Events (race weekends): `/{sport}/events`, `/{sport}/events/{id}`

**Use for:** the schedule, and finding the most recent race.

- Sorted **newest first** by default.
- Page 1 alone is enough to find the most recent completed race. No pagination needed.
- `schedule[]` lists the weekend's sessions (practice, qualifying, race).
- `previousEventId` / `nextEventId` link to the neighboring events.

| Group    | Fields                                                                         |
| -------- | ------------------------------------------------------------------------------ |
| Event    | `id`, `name`, `dateStart`, `dateEnd`, `status`, `sportId`                      |
| Location | `location` (`id`, `name`, `city`, `country`)                                   |
| Season   | `season` (`id`, `name`, `year`)                                                |
| Links    | `previousEventId`, `nextEventId`                                               |
| Sessions | `schedule[]`, each with `id`, `type`, `name`, `status`, `startTime`, `endTime` |

### Session results: `/{sport}/events/{eventId}/sessions/{sessionId}`

**Use for:** race and session results. This is the main results endpoint.

**Key fact:** there's one result row per **(car, individual driver)**.

- A 3-driver WEC car produces 3 rows. They're identical except for `driver` (confirmed on a real WEC race).
- So per-driver endurance scoring works directly, with no car-bundling workaround.
- `position` and `points` are **strings** here, e.g. `"1"` and `"25.0"`.
- The session detail embeds `results[]`. Each row is `V1SessionResultDto` in the spec.
- Several fields are `null` outside "enriched" F1 sessions, so F1 gets the most free-tier depth.

| Group                    | Fields                                                                                                  |
| ------------------------ | ------------------------------------------------------------------------------------------------------- |
| Result                   | `id`, `position` (string), `points` (string), `status` (per sport, see below), `laps`, `gridPosition`   |
| Driver                   | `driver` (`id`, `firstName`, `lastName`, `code`, `number`), `carNumber`                                 |
| Car and team             | `team` (`id`, `name`, `shortName`, `color`), `chassis`, `engineManufacturer`                            |
| Timing                   | `lapTime`, `displayTime`, `gap`, `interval`, `bestLapNumber`, `fastestLap` (`rank`, `time`, `lap`)      |
| Null outside enriched F1 | `pitStops`, `bestLapTime`, `sectors` (`s1`, `s2`, `s3`), `tireStrategy[]`, `q1Time`, `q2Time`, `q3Time` |

### Seasons: `/{sport}/seasons`, `/{sport}/seasons/{id}`

**Use for:** the season list and "round X of Y".

- `roundCount` gives the "Y" in "round X of Y".
- The real `status` value seen is `"Created"`.
- The docs' examples (`"completed"`, `"scheduled"`) don't match reality. Don't hardcode assumed status strings.

### Season drivers: `/{sport}/seasons/{id}/drivers`

**Why it matters:** the best free-tier source for fantasy-scoring depth found so far.

**Confirmed for:** F1 only.

| Group       | Fields                                                                                                                     |
| ----------- | -------------------------------------------------------------------------------------------------------------------------- |
| Races       | `raceStarts`, `raceWins`, `racePodiums`, `topFiveFinishes`, `topTenFinishes`, `bestRaceFinish`, `avgRaceFinishingPosition` |
| Sprints     | `sprintStarts`, `sprintWins`, `sprintPodiums`, `avgSprintFinishingPosition`                                                |
| Qualifying  | `polePositions`, `avgQualifyingPosition`                                                                                   |
| Reliability | `dnfs`, `completionRate`, `totalLapsCompleted`                                                                             |
| Other       | `fastestLaps`, `pointsPerRace`                                                                                             |
| Teams       | `teams[].participationRounds`: the round numbers raced for each team                                                       |

- Gaps in `participationRounds` mean a mid-season team or seat change.

### Standings: `/{sport}/standings/drivers`, `/{sport}/standings/constructors`

**Use for:** championship tables. WEC uses a different endpoint (next section).

**F1/NASCAR fields:** `id`, `position` (number), `points` (string, e.g. `"292.00"`), `firstName`, `lastName`, `code`, `number`, `teams[]`

- `teams[]` is often empty here.
- The same team data shows up populated elsewhere: it was empty on `/standings/drivers` but populated, with `participationRounds`, on `/seasons/{id}`.

### WEC standings: `/wec/standings`

**Key fact:** the "drivers" table is per **crew**, not per driver.

- `name` is one joined string, like `"Fuoco · Molina · Nielsen"`.
- There are no driver IDs.
- To get one WEC driver's season points, sum their session `results` yourself. There's no per-driver WEC standings endpoint.

Standings are split by class. Each table has:

| Field        | Meaning                                                        |
| ------------ | -------------------------------------------------------------- |
| `kind`       | `drivers`, `teams` or `manufacturers`                          |
| `code`       | Table code                                                     |
| `name`       | Table name                                                     |
| `classLabel` | e.g. `"Hypercar"`, `"LMGT3"`; `null` for combined-class titles |
| `rows[]`     | `position`, `name`, `country`, `points`, `isChampion`          |

### NASCAR stages: `/{sport}/events/{eventId}/sessions/{sessionId}/stages`

**Use for:** per-stage results. Free tier.

**Problem:** rows only have a flat `driverName` string, with **no driver ID**.

- Matching rows to our internal `Driver` records needs fragile name/car-number matching.
- Everywhere else, it's a clean ID join.

### Locations: `/{sport}/locations`, `/{sport}/locations/{id}`

**Use for:** circuit and venue data.

## Same field, different types

**Rule:** give every endpoint its own Go struct. Never share a `points` or `position` type across endpoints.

| Field      | Endpoint                                      | Type                                                           |
| ---------- | --------------------------------------------- | -------------------------------------------------------------- |
| `points`   | `/standings/drivers`, session `results`       | **string** (`"25.0"`, `"292.00"`)                              |
| `points`   | `/seasons/{id}`, `/seasons/{id}/drivers`      | **number** (`292`)                                             |
| `position` | `/standings/drivers`, `/seasons/{id}/drivers` | **number**                                                     |
| `position` | session `results`                             | **string** (`"1"`, or `"DNF"`/`"NC"`/`"DSQ"` for unclassified) |

## Status values per sport

**Source:** the spec's description of the `status` field.

**Rule:** branch on the sport. Don't assume one shared enum.

| Series                     | Classified finish                                            | Other values                                                                                           |
| -------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| F1                         | `OK`                                                         | `DNF`, `DNS`, `DSQ`                                                                                    |
| F2, F3, Formula E          | `null`                                                       | `DNF`, `DNS`, `DSQ`                                                                                    |
| NASCAR (Cup/Xfinity/Truck) | `Running`, for nearly every finisher, on the lead lap or not | The retirement reason as free text: `Accident`, `Engine`, `Suspension`, `Brakes`, `DVP`, `Too Slow`, … |
| IndyCar                    | (not stated; `Running` is the likely candidate)              | `Running`, `Contact`, `Mechanical`, `Retired`, `Off Course`, `DNS` (full list as given)                |
| MotoGP/Moto2/Moto3         | `INSTND`                                                     | `OUTSTND` (retired), `NOTFINISHFIRST`, `NOTSTARTED`, `DISQUALIFIED`, `NOTONRESTARTGRID`, `OUTOFLAPS`   |
| WEC                        | `Classified`                                                 | `Retired`, `Not classified`, `Disqualified`                                                            |

## Free vs paid

Confirmed by live testing, not just from the docs.

### Free tier

- `drivers`, `teams`, `events`
- `seasons`, including `/seasons/{id}/drivers`
- `standings`
- session `results`

### Paid plan (HTTP 402)

- **All live endpoints:** `live/sessions/{id}/timing`, `/state`, `/weather`, `/race-control`; NASCAR `live/leaderboard` and `status`; IndyCar live status/timing.
- **All per-lap and telemetry grids:** `lap-times`, `lap-chart`, `telemetry/*`.
- **All IndyCar session sub-resources:** `pit-stops`, `tyres`, `sections`, `flags`, `race-summary`. Every one was confirmed to return 402 on the free tier.

### Commercial plan

- Bulk exports: `/export/results`, `/export/laps`.
- They stream a whole season instead of making paginated calls.

### What a paid-tier rejection looks like

HTTP 402, with the same shape on every endpoint tested:

```json
{
  "message": "This endpoint requires a paid plan",
  "requiredTier": "hobby",
  "currentTier": "free",
  "upgradeUrl": "..."
}
```
