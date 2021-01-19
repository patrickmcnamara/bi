package bi

import (
	"fmt"
	"image"
	"io"
)

// Encode encodes the Image m to the Writer w in BI format using the Model mod.
// This is usually a very lossy process.
func Encode(w io.Writer, m image.Image, mod Model) error {
	hdr := MagicNumber + "," + mod.ID() + "\n"
	if _, err := w.Write([]byte(hdr)); err != nil {
		return err
	}
	mp := m.Bounds().Max
	mx, my := mp.X, mp.Y
	for y := 0; y < my; y++ {
		for x := 0; x < mx; x++ {
			n := mod.ColorToName(m.At(x, y))
			if _, err := fmt.Fprint(w, n); err != nil {
				return err
			}
			if x != mx-1 {
				if _, err := fmt.Fprint(w, ","); err != nil {
					return err
				}
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	return nil
}
