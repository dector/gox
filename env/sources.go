package env

import "os"

type osSource struct{}

func (osSource) LookupEnv(key string) (string, bool) { return os.LookupEnv(key) }

func (m SrcMap) LookupEnv(key string) (value string, found bool) {
	value, found = m[key]
	return value, found
}
