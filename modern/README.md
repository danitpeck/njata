# NJATA Modern Development

## Start Here: Design & Architecture

→ [FRAMEWORK.md](FRAMEWORK.md) — Design philosophy and document overview  
→ [MVP_ROADMAP.md](MVP_ROADMAP.md) — What you're building (Week 1 specification)  
→ [SCHOLAR_STUDY_DESIGN.md](SCHOLAR_STUDY_DESIGN.md) — Core feature details (Scholar Study system)

## Quick Start: Build & Run

Build all packages:

```powershell
go build ./...
```

Run the server:

```powershell
go run ./cmd/njata -port 4000
```

Connect with telnet:

```powershell
telnet localhost 4000
```

If Telnet is not installed, enable it in "Turn Windows features on or off" or use any TCP client.

## Tests

Run all tests:

```powershell
go test ./...
```

## Project Structure

- `cmd/njata/` — Main server executable
- `cmd/test-client/` — Integration test client
- `internal/` — Core packages:
  - `commands/` — Player commands (cast, study, say, etc.)
  - `skills/` — Spell system and proficiency tracking
  - `game/` — Game loop, tick handling
  - `netserver/` — TCP server, connection handling
  - `persist/` — Player save/load
  - `parser/` — Command parser
  - `area/` — Room and area loading
- `skills/`, `classes/`, `races/` — Game data (JSON)
- `players/` — Player save files (generated at runtime)

## Resources

For historical context and legacy code, see [../legacy](../legacy).

For lore and worldbuilding, see [../lore](../lore).
- Commands: look, say, who, help, quit, plus movement (n/s/e/w/u/d and diagonals).
- Legacy SMAUG code and data live under /legacy.
