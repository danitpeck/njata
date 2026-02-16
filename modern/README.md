# Njata Modern (Go)

This is a minimal, playable vertical slice of the Njata MUD rewritten in Go.
It includes a TCP server, basic login, a command loop, and a single in-memory room.

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

- No database or area file loading yet.
- The world is a single in-memory room.
- Commands: look, say, who, help, quit.
- Legacy SMAUG code and data live under /legacy.
