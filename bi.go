package bi

import (
	"compress/gzip"
	"encoding/csv"
	"errors"
	"fmt"
	"image"
	"io"
	"strings"
)

// Encode encodes the Image m to the Writer w in BI format. As BI only supports
// 148 colors, this process is very lossy.
func Encode(w io.Writer, m image.Image) error {
	if _, err := w.Write([]byte("bi\n")); err != nil {
		return err
	}
	return encode(w, m)
}

// EncodeZ encodes the Image m to the Writer w in BIZ format. As BIZ only
// supports 148 colors, this process is very lossy. This is the same as the BI
// format but it is gzipped.
func EncodeZ(w io.Writer, m image.Image) error {
	if _, err := w.Write([]byte("biz\n")); err != nil {
		return err
	}
	gzw := gzip.NewWriter(w)
	if err := encode(gzw, m); err != nil {
		return err
	}
	return gzw.Close()
}

func encode(w io.Writer, m image.Image) error {
	mp := m.Bounds().Max
	mx, my := mp.X, mp.Y
	for y := 0; y < my; y++ {
		for x := 0; x < mx; x++ {
			_, n := NearestColor(m.At(x, y))
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

// Decode reads a BI image from Reader r and returns it as an Image.
func Decode(r io.Reader) (m image.Image, err error) {
	hdr := make([]byte, 3)
	if _, err = r.Read(hdr); err != nil {
		return
	}
	if string(hdr) != "bi\n" {
		err = errors.New("expected header \"bi\\n\" not found")
		return
	}
	return decode(r)
}

// DecodeZ reads a BIZ image from Reader r and returns it as an Image.
func DecodeZ(r io.Reader) (m image.Image, err error) {
	hdr := make([]byte, 4)
	if _, err = r.Read(hdr); err != nil {
		return
	}
	if string(hdr) != "biz\n" {
		err = errors.New("expected header \"biz\\n\" not found")
		return
	}
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return
	}
	return decode(gzr)
}

func decode(r io.Reader) (m image.Image, err error) {
	cs, err := csv.NewReader(r).ReadAll()
	if err != nil || len(cs) < 1 {
		return
	}
	mx := len(cs[0])
	my := len(cs)
	img := image.NewRGBA(image.Rect(0, 0, mx, my))
	for y, row := range cs {
		if len(row) != mx {
			err = fmt.Errorf("row %d should be %d long, not %d", y, mx, len(row))
			return
		}
		for x, n := range row {
			c := colors[strings.TrimSpace(n)]
			img.SetRGBA(x, y, c)
		}
	}
	m = img
	return
}

func init() {
	image.RegisterFormat("bi", "bi\n", Decode, func(io.Reader) (c image.Config, err error) { return })
	image.RegisterFormat("biz", "biz\n", DecodeZ, func(io.Reader) (c image.Config, err error) { return })
}
