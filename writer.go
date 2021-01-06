package bi

import (
	"fmt"
	"image"
	"io"
)

// Encode encodes the Image m to the Writer w in BI format. As BI only supports
// 148 colors, this process is very lossy.
func Encode(w io.Writer, m image.Image) error {
	if _, err := w.Write([]byte("bi\n")); err != nil {
		return err
	}
	mp := m.Bounds().Max
	mx, my := mp.X, mp.Y
	for y := 0; y < my; y++ {
		for x := 0; x < mx; x++ {
			_, n := nearestColor(m.At(x, y))
			if _, err := fmt.Fprint(w, n); err != nil {
				return err
			}
			if x != mx-1 {
				if _, err := fmt.Fprint(w, ", "); err != nil {
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
