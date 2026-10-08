package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	_ Source = SrcMap(nil)
	_ Source = osSource{}
)

func TestSrcMapLookupEnv(t *testing.T) {
	for _, tc := range []struct {
		name  string
		src   SrcMap
		want  string
		found bool
	}{
		{"nil", nil, "", false},
		{"unset", SrcMap{}, "", false},
		{"empty", SrcMap{"KEY": ""}, "", true},
		{"value", SrcMap{"KEY": "value"}, "value", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			is := assert.New(t)

			value, found := tc.src.LookupEnv("KEY")
			is.Equal(tc.want, value)
			is.Equal(tc.found, found)
		})
	}
}

func TestGet(t *testing.T) {
	is := assert.New(t)

	for _, src := range []SrcMap{nil, {}, {"KEY": ""}, {"KEY": "value"}} {
		want, wantFound := src["KEY"]
		got, found := Get(src, "KEY")
		is.Equal(want, got)
		is.Equal(wantFound, found)
	}
}

func TestSrcOSLookupEnv(t *testing.T) {
	is := assert.New(t)

	const key = "TOW_PKG_ENV_TEST_LOOKUP"
	// Setenv records the original value for cleanup before testing unset.
	t.Setenv(key, "initial")
	is.NoError(os.Unsetenv(key))
	src := SrcOS()
	value, found := Get(src, key)
	is.Empty(value)
	is.False(found)
	for _, want := range []string{"", "value"} {
		t.Setenv(key, want)
		value, found := Get(src, key)
		is.Equal(want, value)
		is.True(found)
	}
}
