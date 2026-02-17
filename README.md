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

Connect:

```powershell
telnet localhost 4000
```

See [modern/README.md](modern/README.md) for details and tests.

## Legacy (C)

The original SMAUG-based code and data live under [legacy/](legacy/).
