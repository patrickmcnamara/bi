package bi

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"image"
	"io"
)

// Decode reads a BI image from Reader r and returns it as an Image m.
func Decode(r io.Reader) (m image.Image, err error) {
	hdrr, err := headerReader(r)
	if err != nil {
		return
	}
	return decode(hdrr)
}

// DecodeConfig decodes the Model and dimensions of a BI image from Reader r.
func DecodeConfig(r io.Reader) (c image.Config, err error) {
	hdrr, err := headerReader(r)
	if err != nil {
		return
	}
	return decodeConfig(hdrr)
}

func headerReader(r1 io.Reader) (r2 io.Reader, err error) {
	bufr := bufio.NewReaderSize(r1, 128)
	if err = decodeHeader(bufr); err != nil {
		return
	}
	r2 = bufr
	return
}

func decodeHeader(r *bufio.Reader) error {
	hdr, err := r.ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("bi: error decoding: %w", err)
	}
	if string(hdr) != MagicNumber {
		return fmt.Errorf("bi: expected magic number %q not found", MagicNumber)
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
			c, ok := CSSColModLvl4.NameToColor(n)
			if !ok {
				err = fmt.Errorf("bi: unexpected colour %q on line %d", n, y+2)
				return
			}
			img.Set(x, y, c)
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
	cfg.ColorModel = CSSColModLvl4
	return
}

func init() {
	image.RegisterFormat("bi", MagicNumber, Decode, DecodeConfig)
}
