package render

import (
	"fmt"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func TestBinSearch(t *testing.T) {
	ttf, err := opentype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("cannot parse Go Regular font: %v", err)
	}
	fontFn := func(size float64) (font.Face, error) {
		face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
			Size: size,
			DPI:  72,
		})
		return face, err
	}
	run := func(maxHeight int, wantAscent, wantDescent fixed.Int26_6) {
		t.Run(fmt.Sprintf("%d", maxHeight), func(t *testing.T) {
			face, err := binsearch(maxHeight, fontFn)
			if err != nil {
				t.Fatalf("binsearch: %v", err)
			}
			// point size is lost (cannot be retrieved from *opentype.Face),
			// so compare other font metrics to make sure we picked the right
			// one.
			if g, e := face.Metrics().Ascent, wantAscent; g != e {
				t.Errorf("wrong ascent: %v != %v", g, e)
			}
			if g, e := face.Metrics().Descent, wantDescent; g != e {
				t.Errorf("wrong descent: %v != %v", g, e)
			}
		})
	}

	run(10, fixed.I(8)|11, fixed.I(1)|53)
	run(20, fixed.I(16)|22, fixed.I(3)|42)
	run(30, fixed.I(24)|33, fixed.I(5)|30)
	run(40, fixed.I(32)|45, fixed.I(7)|19)
	run(50, fixed.I(40)|56, fixed.I(9)|8)
	run(60, fixed.I(49)|3, fixed.I(10)|61)
}
