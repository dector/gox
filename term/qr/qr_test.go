package qr

import (
	"bytes"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/boombuler/barcode/qr"
	"github.com/stretchr/testify/assert"
)

var updateGolden = flag.Bool("update", false, "update golden files")

func TestWriteTerminal(t *testing.T) {
	for _, tc := range []struct {
		name string
		text string
	}{
		{"hello", "hello"},
		{"uri", "https://dector.space/gox/term/qr"},
		{"utf8", "自由"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			is := assert.New(t)

			var out bytes.Buffer
			if !is.NoError(WriteTerminal(&out, tc.text)) {
				return
			}

			assertGolden(t, tc.name, out.String())
			assertTerminalModules(t, tc.text, out.String())
		})
	}
}

func assertGolden(t *testing.T, name, got string) {
	t.Helper()
	is := assert.New(t)

	goldenPath := filepath.Join("testdata", name+".golden")
	if *updateGolden {
		if !is.NoError(os.MkdirAll("testdata", 0o755)) {
			return
		}
		if !is.NoError(os.WriteFile(goldenPath, []byte(got), 0o644)) {
			return
		}
	}
	want, err := os.ReadFile(goldenPath)
	if !is.NoError(err) {
		return
	}
	is.Equal(string(want), got, "golden file %s", goldenPath)
}

func assertTerminalModules(t *testing.T, text, output string) {
	t.Helper()
	is := assert.New(t)

	image, err := qr.Encode(text, qr.M, qr.Auto)
	if !is.NoError(err) {
		return
	}
	const (
		quietZone  = 4 // White border around the QR code, in modules.
		colorStart = "\x1b[30;47m"
		colorReset = "\x1b[0m"
	)

	size := image.Bounds().Dx()
	width := size + 2*quietZone
	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")
	is.Len(lines, (width+1)/2, "each terminal row holds two module rows")
	for row, line := range lines {
		is.True(strings.HasPrefix(line, colorStart), "row %d missing ANSI colors", row)
		is.True(strings.HasSuffix(line, colorReset), "row %d missing ANSI reset", row)
		cells := []rune(strings.TrimSuffix(strings.TrimPrefix(line, colorStart), colorReset))
		is.Len(cells, width, "columns in row %d", row)
		for col, cell := range cells {
			is.Contains(" █▀▄", string(cell), "unexpected cell %q", cell)
			// Compare the upper and lower halves of each terminal cell.
			for half := range 2 {
				x, y := col-quietZone, row*2+half-quietZone
				wantDark := false
				if x >= 0 && y >= 0 && x < size && y < size {
					r, g, b, _ := image.At(x, y).RGBA()
					wantDark = r+g+b < 3*32768
				}
				gotDark := cell == '█' || (half == 0 && cell == '▀') || (half == 1 && cell == '▄')
				is.Equal(wantDark, gotDark, "module (%d,%d): dark", x, y)
			}
		}
	}
}

func TestWriteTerminalEncodingError(t *testing.T) {
	is := assert.New(t)

	var out bytes.Buffer
	err := WriteTerminal(&out, strings.Repeat("secret", 10000))
	if is.Error(err) {
		is.NotContains(err.Error(), "secret", "encoding error exposed the payload")
	}
	is.Zero(out.Len(), "encoding failure wrote output")
}

type failingWriter struct {
	err   error
	calls int
}

func (w *failingWriter) Write(p []byte) (int, error) {
	w.calls++
	return 0, w.err
}

func TestWriteTerminalWriterError(t *testing.T) {
	is := assert.New(t)

	want := errors.New("write failed")
	out := &failingWriter{err: want}
	is.ErrorIs(WriteTerminal(out, "hello"), want)
	is.Equal(1, out.calls, "continued after write failure")
}
