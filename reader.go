package bi

import (
	"encoding/csv"
	"fmt"
	"image"
	"io"
	"strings"
)

// Decode reads a BI image from Reader r and returns it as an Image m.
func Decode(r io.Reader) (m image.Image, err error) {
	if err = decodeMagic(r); err != nil {
		return
	}
	return decode(r)
}

// DecodeConfig decodes the Model and dimensions of a BI image from Reader r.
func DecodeConfig(r io.Reader) (c image.Config, err error) {
	if err = decodeMagic(r); err != nil {
		return
	}
	return decodeConfig(r)
}

func decodeMagic(r io.Reader) error {
	hdr := make([]byte, len(Magic))
	if _, err := r.Read(hdr); err != nil {
		return image.ErrFormat
	}
	if string(hdr) != Magic {
		return fmt.Errorf("bi: expected magic number %q not found", Magic)
	}
	return nil
}

func decode(r io.Reader) (m image.Image, err error) {
	cs, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return
	}
	mx := len(cs[0])
	my := len(cs)
	img := image.NewRGBA(image.Rect(0, 0, mx, my))
	for y, row := range cs {
		if len(row) != mx {
			err = fmt.Errorf("bi: row %d should have length %d, was %d", y, mx, len(row))
			return
		}
		for x, n := range row {
			c, ok := colors[strings.TrimSpace(n)]
			if !ok {
				err = image.ErrFormat
				return
			}
			img.SetRGBA(x, y, c)
		}
	}
	m = img
	return
}

func decodeConfig(r io.Reader) (cfg image.Config, err error) {
	cs, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return
	}
	cfg.Width = len(cs[0])
	cfg.Height = len(cs)
	cfg.ColorModel = CSSColModLevel4
	return
}

func init() {
	image.RegisterFormat("bi", "bi\n", Decode, DecodeConfig)
}
