# Njata

This repo contains a modern Go rewrite of the Njata MUD alongside the legacy SMAUG-based C code.

## Folder Structure (Top Level)

- areas/ - Game world area data in JSON (rooms, mobiles, objects).
- backup/ - Archived or migrated content (including legacy area backups).
- cmd/ - Entry points for binaries (njata server).
- internal/ - Core Go server code (game loop, commands, networking, loaders).
- legacy/ - Original SMAUG C codebase and data kept for reference.
- lore/ - Worldbuilding references, narrative notes, and design docs.
- modern/ - Project docs, roadmap, and MVP notes for the Go rewrite.
- players/ - Player save files (JSON).
- races/ - Race definitions and data.
- scripts/ - Utility scripts and test helpers.
- skills/ - Skills/abilities data and definitions.
- system/ - Misc system data and configuration.

## Modern (Go)

Build and run from the repo root on Windows 10/11:

```powershell
go build ./...
```

```powershell
go run ./cmd/njata -port 4000
```

What this means:

- go run = builds a temporary copy and starts it right away (nothing new appears in your folder).
- go build = creates a real njata.exe you can run later.

Build a binary (creates njata.exe in the repo root):

```powershell
go build ./cmd/njata
```

Run the compiled binary:

```powershell
./njata.exe -port 4000
```

Connect:

```powershell
telnet localhost 4000
```

## Tests

Unit tests (Go):

```powershell
go test ./...
```

Integration scripts (Python) live in scripts/tests and expect the server running on port 4000:

```powershell
go run ./cmd/njata -port 4000
```

Then in another terminal:

```powershell
python scripts/tests/test_spawn.py
python scripts/tests/test_spell_combat.py
```

See [modern/README.md](modern/README.md) for details and tests.

## Legacy (C)

The original SMAUG-based code and data live under [legacy/](legacy/).
