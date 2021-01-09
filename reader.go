package bi

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"image"
	"io"
	"strings"
)

// Decode reads a BI image from Reader r and returns it as an Image m.
func Decode(r io.Reader) (m image.Image, err error) {
	hdrr, mod, err := headerReader(r)
	if err != nil {
		return
	}
	return decode(hdrr, mod)
}

// DecodeConfig decodes the Model and dimensions of a BI image from Reader r.
func DecodeConfig(r io.Reader) (c image.Config, err error) {
	hdrr, mod, err := headerReader(r)
	if err != nil {
		return
	}
	return decodeConfig(hdrr, mod)
}

func headerReader(r1 io.Reader) (r2 io.Reader, mod Model, err error) {
	bufr := bufio.NewReaderSize(r1, 128)
	mod, err = decodeHeader(bufr)
	if err != nil {
		return
	}
	r2 = bufr
	return
}

func decodeHeader(r *bufio.Reader) (mod Model, err error) {
	hdr, err := r.ReadBytes('\n')
	if err != nil {
		err = fmt.Errorf("bi: error decoding: %w", err)
		return
	}
	tkns := strings.Split(strings.TrimSuffix(string(hdr), "\n"), ",")
	if len(tkns) < 1 || tkns[0] != MagicNumber {
		err = fmt.Errorf("bi: expected magic number %q not found", MagicNumber)
		return
	}
	if len(tkns) < 2 || tkns[1] == "" {
		err = fmt.Errorf("bi: expected color model name not found")
		return
	}
	if len(tkns) > 2 {
		err = fmt.Errorf("bi: invalid header, too many tokens")
		return
	}
	val, ok := models.Load(tkns[1])
	if !ok {
		err = fmt.Errorf("bi: color model %q is invalid", tkns[1])
	}
	mod, _ = val.(Model)
	return
}

func decode(r io.Reader, mod Model) (m image.Image, err error) {
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
			c, ok := mod.NameToColor(n)
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

func decodeConfig(r io.Reader, mod Model) (cfg image.Config, err error) {
	cs, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return
	}
	cfg.Width = len(cs[0])
	cfg.Height = len(cs)
	cfg.ColorModel = mod
	return
}

func init() {
	image.RegisterFormat("bi", MagicNumber, Decode, DecodeConfig)
}
