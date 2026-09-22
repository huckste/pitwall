# OCB API notes

Detailed endpoint findings for Orange Cat Blacktop (`https://api.ocblacktop.com/v1`). Auth: `x-api-key` header. Full OpenAPI spec (not checked into this repo): `~/Downloads/api-1.json`.

Free tier: 7,500 requests/month, 60/minute. Confirmed paid-tier rejection shape (HTTP 402), consistent across every endpoint tested:
```json
{"message": "This endpoint requires a paid plan", "requiredTier": "hobby", "currentTier": "free", "upgradeUrl": "..."}
```

## Series covered

F1, F2, F3, Formula E, IndyCar, NASCAR (Cup/Truck/Xfinity), MotoGP/Moto2/Moto3, WEC, WRC. **No IMSA** — not available anywhere in this API.

## Endpoint shape per series (pattern: `/{sport}/...`)

`{sport}` in the URL path is the lowercase, dash-for-compound slug (`formula1`, `nascar`, `nascar-trucks`, `indycar`, `wec`, ...). This is DIFFERENT from the `sportId` field that comes back inside response bodies, which is always the dashed form (`"formula-1"`) even for series whose URL slug has no dash. Don't conflate the two.

- `/{sport}/drivers`, `/{sport}/drivers/{id}` — driver list/detail. F1, Formula E, F2 all confirmed identical shape: `id, firstName, lastName, birthDate, number, tla, country{name,twoCode,threeCode}`.
- `/{sport}/teams`, `/{sport}/teams/{id}` — team detail includes a `drivers[]` roster array (each with `teamSeasonParticipations[]`) and `seasonParticipations[]`. Supports "follow by manufacturer."
- `/{sport}/events`, `/{sport}/events/{id}` — event (race weekend) list/detail. Fields: `id, name, dateStart, dateEnd, status, location{id,name,city,country}, sportId, season{id,name,year}, previousEventId, nextEventId, schedule[]` (schedule = practice/qualifying/race sessions with `id,type,name,status,startTime,endTime`). Sorted newest-first by default; page 1 alone is enough to find "most recent completed race" without pagination.
- `/{sport}/events/{eventId}/sessions/{sessionId}` — session detail + `results[]` embedded. This is the main results endpoint. Each result row (`V1SessionResultDto` in the spec) is **one row per (car, individual driver)** — confirmed empirically via a real WEC race: a 3-driver car produces 3 rows sharing identical `position`, `carNumber`, `team`, `points`, `laps`, `status`, only `driver` differs. This means real per-individual-driver endurance scoring is directly buildable, no car-bundling workaround needed.
  - Fields: `id, position (STRING, e.g. "1"), lapTime, displayTime, laps, driver{id,firstName,lastName,code,number}, carNumber, team{id,name,shortName,color}, chassis, engineManufacturer, fastestLap{rank,time,lap}, status (STRING, per-sport vocab — see below), points (STRING, e.g. "25.0"), gap, interval, pitStops, bestLapTime, bestLapNumber, sectors{s1,s2,s3}, tireStrategy[], gridPosition, q1Time, q2Time, q3Time`.
  - Many fields (`pitStops`, `bestLapTime`, `sectors`, `tireStrategy`, `q1-3Time`) are null outside "enriched F1 sessions" — F1 gets more free-tier depth than other series.
- `/{sport}/seasons`, `/{sport}/seasons/{id}` — season list/detail. `roundCount` field solves "round X of Y" cleanly. Real `status` value seen: `"Created"` (NOT `"completed"`/`"scheduled"` as the docs' example implies — don't hardcode assumed status strings).
- `/{sport}/seasons/{id}/drivers` — **rich free-tier season-aggregate stats**: `raceStarts, raceWins, racePodiums, avgRaceFinishingPosition, topFiveFinishes, topTenFinishes, bestRaceFinish, sprintStarts/Wins/Podiums, avgSprintFinishingPosition, polePositions, avgQualifyingPosition, fastestLaps, dnfs, completionRate, totalLapsCompleted, pointsPerRace`, plus `teams[].participationRounds` (array of round numbers, gaps = mid-season team/seat changes). Confirmed for F1; this is the best free-tier source for fantasy-scoring depth found so far.
- `/{sport}/standings/drivers`, `/{sport}/standings/constructors` (or `/{sport}/standings` for WEC, class-scoped) — championship standings.
  - F1/NASCAR shape: `id, position (NUMBER), points (STRING, e.g. "292.00"), firstName, lastName, code, number, teams[]` (teams array often empty even when populated elsewhere — seen empty on `/standings/drivers` but populated with `participationRounds` on `/seasons/{id}`).
  - **WEC `/wec/standings` "drivers" table is CREW-level, not per-driver**: `name` is a single joined string like `"Fuoco · Molina · Nielsen"`, no driver IDs at all. To get an individual WEC driver's season points, sum their session `results` yourself — there is no per-driver WEC standings endpoint.
  - WEC standings are class-scoped: response has `kind` (`drivers`/`teams`/`manufacturers`), `code`, `name`, `classLabel` (e.g. `"Hypercar"`, `"LMGT3"`, null for combined-class titles), `rows[]` (`position, name, country, points, isChampion`).
- `/{sport}/locations`, `/{sport}/locations/{id}` — circuit/venue data.

## Type inconsistency warning (design every Go struct per-endpoint, don't share types)

| Field | Endpoint | Type |
|---|---|---|
| `points` | `/standings/drivers`, session `results` | **string** (`"25.0"`, `"292.00"`) |
| `points` | `/seasons/{id}`, `/seasons/{id}/drivers` | **number** (`292`) |
| `position` | `/standings/drivers`, `/seasons/{id}/drivers` | **number** |
| `position` | session `results` | **string** (`"1"`, or `"DNF"`/`"NC"`/`"DSQ"` for unclassified) |

## Status vocabulary (per-sport, from the spec description — branch on it, don't assume one enum)

- **F1**: `OK` (classified finish) else `DNF`, `DNS`, `DSQ`.
- **F2, F3, Formula E**: `null` for classified finish, else `DNF`, `DNS`, `DSQ`.
- **NASCAR** (Cup/Xfinity/Truck): `Running` for nearly every finisher (on lead lap or not); anything else is the retirement reason as free text (`Accident`, `Engine`, `Suspension`, `Brakes`, `DVP`, `Too Slow`, ...).
- **IndyCar**: `Running`, `Contact`, `Mechanical`, `Retired`, `Off Course`, `DNS`.
- **MotoGP/Moto2/Moto3**: `INSTND` (classified), `OUTSTND` (retired), `NOTFINISHFIRST`, `NOTSTARTED`, `DISQUALIFIED`, `NOTONRESTARTGRID`, `OUTOFLAPS`.
- **WEC**: `Classified`, `Retired`, `Not classified`, `Disqualified`.

## Paid-tier gating (confirmed by live testing, not just docs)

- **Free tier**: `drivers`, `teams`, `events`, `seasons` (including `/seasons/{id}/drivers`), `standings`, session `results`.
- **Paid plan required** (HTTP 402): all live endpoints (`live/sessions/{id}/timing`, `/state`, `/weather`, `/race-control`, NASCAR `live/leaderboard`/`status`, IndyCar live status/timing), all per-lap/telemetry grids (`lap-times`, `lap-chart`, `telemetry/*`), and **all of IndyCar's session sub-resources** (`pit-stops`, `tyres`, `sections`, `flags`, `race-summary` — confirmed every one 402s on free tier).
- **Commercial plan required**: bulk `/export/results`, `/export/laps` (stream a whole season instead of paginated calls).

## NASCAR-specific note

`/{sport}/events/{eventId}/sessions/{sessionId}/stages` (free tier) gives per-stage results, but rows only have a flat `driverName` string — **no driver ID**. Joining to our internal `Driver` records will need fragile name/car-number matching here, not a clean ID join like everywhere else.
