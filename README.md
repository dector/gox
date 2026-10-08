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

## Development

Run from the repository root:

```sh
go test ./env/...
go vet ./env/...
```

To verify `env` independently of the workspace:

```sh
cd env
GOWORK=off go test ./...
```

### Adding new modules

Add future library modules with `go work use ./<library>`.
Release `env` with tags such as `env/v0.1.0`.
