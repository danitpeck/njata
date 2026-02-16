# Njata Modern (Go)

This is a minimal, playable vertical slice of the Njata MUD rewritten in Go.
It includes a TCP server, basic login, a command loop, area room loading, and basic player persistence.

## Requirements

- Windows 10/11
- Go 1.26 (latest stable) (https://go.dev/dl/)

## Build

```powershell
go build ./...
```

## Run

```powershell
go run ./cmd/njata -port 4000
```

## Connect

```powershell
telnet localhost 4000
```

If Telnet is not installed, enable it in "Turn Windows features on or off" or use any TCP client.

## Tests

```powershell
go test ./...
```

## Notes

- Rooms are loaded from the .are files in /areas (room section only).
- Player data is stored as JSON under /players.
- Commands: look, say, who, help, quit, plus movement (n/s/e/w/u/d and diagonals).
- Legacy SMAUG code and data live under /legacy.
