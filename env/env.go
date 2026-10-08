// Package env provides environment lookup sources without application-specific
// defaults or validation. Callers should resolve configuration at startup and
// inject the resulting snapshot rather than reading the environment per request.
package env

// Source resolves a value by environment variable name. A missing variable
// reports found=false; an explicitly empty variable reports found=true.
type Source interface {
	LookupEnv(key string) (value string, found bool)
}

// Get returns a value and whether the variable was found. An explicitly empty
// variable returns "", true; an unset variable returns "", false.
func Get(src Source, key string) (value string, found bool) {
	return src.LookupEnv(key)
}

// SrcOS returns a source backed by the live process environment.
func SrcOS() Source { return osSource{} }

// SrcMap is an in-memory source for tests and explicit configuration.
type SrcMap map[string]string
