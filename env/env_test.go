package env

import (
	"os"
	"testing"
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
			value, found := tc.src.LookupEnv("KEY")
			if value != tc.want || found != tc.found {
				t.Fatalf("LookupEnv = %q, %v; want %q, %v", value, found, tc.want, tc.found)
			}
		})
	}
}

func TestGet(t *testing.T) {
	for _, src := range []SrcMap{nil, {}, {"KEY": ""}, {"KEY": "value"}} {
		want, wantFound := src["KEY"]
		if got, found := Get(src, "KEY"); got != want || found != wantFound {
			t.Fatalf("Get = %q, %v; want %q, %v", got, found, want, wantFound)
		}
	}
}

func TestSrcOSLookupEnv(t *testing.T) {
	const key = "TOW_PKG_ENV_TEST_LOOKUP"
	// Setenv records the original value for cleanup before testing unset.
	t.Setenv(key, "initial")
	if err := os.Unsetenv(key); err != nil {
		t.Fatal(err)
	}
	src := SrcOS()
	if value, found := Get(src, key); value != "" || found {
		t.Fatalf("unset = %q, %v", value, found)
	}
	for _, want := range []string{"", "value"} {
		t.Setenv(key, want)
		if value, found := Get(src, key); value != want || !found {
			t.Fatalf("LookupEnv = %q, %v; want %q, true", value, found, want)
		}
	}
}
