# gox

GO eXtensions. A monorepo of small, independently versioned Go libraries.

## env

Environment variables lookup.

```go
import "dector.space/gox/env"

src := env.SrcOS()
value, found := env.Get(src, "PORT")

// Supply explicit configuration in tests.
testSource := env.SrcMap{"DATABASE_PATH": "test.db"}
```

`SrcOS()` reads the live process environment. Resolve application configuration at
startup and inject the resulting snapshot rather than reading it per request.
`Get` returns `(value, found)`: an unset variable returns `"", false`, while
an explicitly empty variable returns `"", true`.

## term/qr

Render QR codes in terminals with UTF-8 half-blocks, explicit ANSI colors, and a
four-module quiet zone.

```go
import "dector.space/gox/term/qr"

err := qr.WriteTerminal(os.Stdout, "https://dector.space")
```

`WriteTerminal` uses medium error correction and returns encoding or writer errors.
The terminal must support UTF-8 and ANSI colors.

## Development

Run from the repository root:

```sh
go test ./env/... ./term/qr/...
go vet ./env/... ./term/qr/...
```

To verify `env` independently of the workspace:

```sh
cd env
GOWORK=off go test ./...
```

### Adding new modules

Add future library modules with `go work use ./<library>`.
Release modules with tags such as `env/v0.1.0` or `term/qr/v0.1.0`.
