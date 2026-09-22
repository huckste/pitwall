# pitwall

## Goal

A dashboard to follow race car drivers (and manufacturers/teams) across multiple racing series, and see their most recent race result at a glance. Low barrier to entry: minimal info at a glance (name, last result, in contention or not), full depth only on click-through — the user gets overwhelmed by too much info at once, so this is a hard design constraint, not a nice-to-have.

Series of interest: NASCAR, F1, WEC, IMSA, F2, F3, IndyCar. IMSA has no usable public data source anywhere found — excluded from v1, revisit only if a data source appears.

Secondary, later phase: a cross-series fantasy league (pick drivers/cars across Formula-type series and Endurance-type series, public/private groups, no real-money/gambling handled by the app). Same app, same backend/auth foundation as the dashboard — fantasy is an opt-in module, not a separate product. Don't build this until the dashboard MVP's data layer is solid.

Context on why this is being built: user is currently unemployed, building this primarily as a portfolio/interview piece. Build with real intent (not a throwaway demo) since the cross-series angle is a genuine gap official single-series apps structurally can't fill — but don't plan around it making money; the real open risk is distribution (no marketing budget), not idea quality.

## Stack

- **Backend**: Go (`backend/`). Idiomatic use: a poller service (goroutines) hits the OCB API on a schedule and writes into Postgres; the HTTP API serves reads from Postgres, not live from OCB. Users interact with our copy of the data, not the live API directly — this bounds API cost regardless of user count.
- **Shared models**: C# class library (`shared/Pitwall.Shared`) — DTOs shared between the two C# frontends.
- **Web**: Blazor WebAssembly (`web/Pitwall.Web`) — calls the Go backend's REST API over HTTP.
- **Mobile**: .NET MAUI (`mobile/`) — same REST API. **MAUI cannot be developed on Linux at all** (the `maui` SDK workload is Windows/macOS-only, confirmed by a failed `dotnet workload install maui` on this machine). User dual-boots Windows on the same machine for MAUI work. iOS builds specifically also need a Mac (Xcode) — planned via a cloud Mac CI service (Codemagic/GitHub Actions macOS runners/MacStadium) when it's time to actually build/test on iOS, not set up yet.
- **Database**: Postgres (chosen over SQL Server despite user's SQL Server background, since SQL Server isn't a natural fit on Linux/in the Go ecosystem; core SQL knowledge transfers directly).
- No C in the stack unless a future live-timing/telemetry decode task turns up a real fit — hasn't so far.

## Data source: Orange Cat Blacktop (OCB)

Base URL `https://api.ocblacktop.com/v1`, auth via `x-api-key` header. Free tier: 7,500 requests/month, 60/minute. Full downloaded OpenAPI spec at `~/Downloads/api-1.json` (not in this repo) — use `jq` against `.paths`/`.components.schemas` for documented shapes, but **cross-check against live responses**, since docs and reality have already diverged more than once.

**See `docs/ocb-api-notes.md` for detailed endpoint-by-endpoint findings** (field shapes, gotchas, tier gating).

Series OCB covers: F1, F2, F3, Formula E, IndyCar, NASCAR (Cup/Truck/Xfinity), MotoGP/Moto2/Moto3, WEC, WRC. No IMSA.

## Design principles

- Two views, not one dense page: glance (name, last result, in contention) + drill-down (full stats) only on demand.
- "Follow" works at the driver level AND the team/manufacturer level (e.g. follow "Porsche" in WEC without tracking which 2-3 rotating drivers are in the car). This also sidesteps endurance racing's multi-driver-per-car attribution problem for the dashboard.
- One app, not two, for dashboard + fantasy — fantasy is additive, not a separate codebase/foundation.
- Testing direction (not yet built): headless automation that drives the REAL rendered UI (real button clicks, not just API calls) via Playwright for the Blazor web app, once there's an actual frontend to drive. For fantasy balance-testing: cheap heuristic bots simulating many seasons/strategies for numeric exploit-finding, LLM agents reserved for a later qualitative "does this feel confusing" pass.

## Known gotchas

- **OCB field types are inconsistent across endpoints for the same field name.** `points` is a STRING in `/standings/drivers` and session `results`, but a NUMBER in `/seasons/{id}/drivers`. `position` is the reverse. Every endpoint needs its own typed Go struct — do not share a `Points`/`Position` type across endpoints.
- **`sportId` in response bodies uses a dashed form** (`"formula-1"`), different from the **URL path slug** (`"formula1"`, no dash) used to call the API. Don't conflate the two.
- **Status vocabulary differs per sport** and doesn't match the docs' example values (e.g. real season `status` is `"Created"`, not `"completed"`/`"scheduled"` as shown in docs examples). Needs a real per-sport normalization layer, not a shared enum assumption.
- **WEC endurance results ARE per-individual-driver**, confirmed empirically: a 3-driver car's session result produces 3 rows, one per co-driver, sharing identical position/points/car/team. Real per-driver endurance fantasy scoring is directly buildable — resolved what was originally assumed to be a hard design problem.
- **But WEC's `/wec/standings` "drivers" table is crew-level, not per-driver** — `name` is a joined string like `"Fuoco · Molina · Nielsen"`, no driver IDs. A WEC driver's season points must be computed ourselves by summing their session results, no standings shortcut exists.
- **Paid-tier gating**: all live endpoints (live timing/standings/status) and all per-lap/telemetry grids are paid-plan (403/402, consistent rejection shape `{message, requiredTier, currentTier, upgradeUrl}`). Bulk `/export/*` dumps are Commercial-plan only. IndyCar's session sub-resources (pit-stops, tyres, sections, flags, race-summary) are ALL paid-gated — don't design free-tier IndyCar features around them. `/formula1/seasons/{id}/drivers` is a notable free-tier exception with rich season-aggregate stats (wins, podiums, avg finish, DNFs, completion rate, etc.) — good fantasy-scoring depth without paid telemetry, at least for F1.
- **NASCAR's `stages` endpoint has no driver ID**, only a flat `driverName` string — joining to our internal `Driver` requires fragile name/car-number matching, not a clean ID join.
- The C# `Series` enum (`shared/Pitwall.Shared/Series.cs`) and the Go `Series` const block (`backend/models.go`) drifted out of sync (Go has 13 values matching OCB's lowercase slugs, C# still has an older 6-value set) — needs reconciling before the frontend actually talks to the backend. Also: C# enum member names can't contain hyphens (OCB slugs like `nascar-trucks` do), and `System.Text.Json` serializes C# enums as integers by default — both need handling (custom converter/attribute) when wiring this up.

## Working style

User is hands-on learning Go/Blazor/MAUI/Postgres by building this themselves. When guiding setup/implementation: give ONE step, then wait for confirmation before the next — don't bundle multiple steps into one response. Trust a "done" reply as complete; only verify actual state (filesystem/git/build) when something seems off or the user asks, not as a standing habit.
