// Package qr renders QR codes for terminals.
package qr

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/boombuler/barcode/qr"
)

// WriteTerminal writes a QR code for text using medium error correction.
// Explicit black/white ANSI colors keep it scannable on either terminal theme.
// Half-block cells render two modules per row, with a four-module quiet zone.
// The writer must support UTF-8 and ANSI color escape sequences.
func WriteTerminal(out io.Writer, text string) error {
	image, err := qr.Encode(text, qr.M, qr.Auto)
	if err != nil {
		// Encoder errors can contain the payload, which may be a secret.
		return errors.New("qr: encoding failed")
	}
	bounds := image.Bounds()
	dark := func(x, y int) bool {
		if x < 0 || y < 0 || x >= bounds.Dx() || y >= bounds.Dy() {
			return false
		}
		r, g, b, _ := image.At(x, y).RGBA()
		return r+g+b < 3*32768
	}
	for y := -4; y < bounds.Dy()+4; y += 2 {
		var line strings.Builder

		line.WriteString("\x1b[30;47m")
		for x := -4; x < bounds.Dx()+4; x++ {
			top, bottom := dark(x, y), dark(x, y+1)
			switch {
			case top && bottom:
				line.WriteString("█")
			case top:
				line.WriteString("▀")
			case bottom:
				line.WriteString("▄")
			default:
				line.WriteString(" ")
			}
		}
		if _, err = fmt.Fprintln(out, line.String()+"\x1b[0m"); err != nil {
			return err
		}
	}
	return nil
}
