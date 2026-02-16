# Njata

This repo contains a modern Go rewrite of the Njata MUD alongside the legacy SMAUG-based C code.

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
