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

	check := func(t *testing.T, lines int, maxHeight int, wantAscent, wantDescent fixed.Int26_6) {
		face, err := binsearch(lines, maxHeight, fontFn)
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
	}

	t.Run("lines=1", func(t *testing.T) {
		run := func(maxHeight int, wantAscent, wantDescent fixed.Int26_6) {
			t.Run(fmt.Sprintf("height=%d", maxHeight), func(t *testing.T) {
				check(t, 1, maxHeight, wantAscent, wantDescent)
			})
		}
		run(10, fixed.I(8)|11, fixed.I(1)|53)
		run(20, fixed.I(16)|22, fixed.I(3)|42)
		run(30, fixed.I(24)|33, fixed.I(5)|30)
		run(40, fixed.I(32)|45, fixed.I(7)|19)
		run(50, fixed.I(40)|56, fixed.I(9)|8)
		run(60, fixed.I(49)|3, fixed.I(10)|61)
	})

	t.Run("lines=2", func(t *testing.T) {
		run := func(maxHeight int, wantAscent, wantDescent fixed.Int26_6) {
			t.Run(fmt.Sprintf("height=%d", maxHeight), func(t *testing.T) {
				check(t, 2, maxHeight, wantAscent, wantDescent)
			})
		}
		run(20, fixed.I(8)|11, fixed.I(1)|53)
		run(30, fixed.I(12)|17, fixed.I(2)|47)
		run(40, fixed.I(16)|22, fixed.I(3)|42)
		run(50, fixed.I(20)|28, fixed.I(4)|36)
		run(60, fixed.I(24)|33, fixed.I(5)|30)
		run(70, fixed.I(28)|39, fixed.I(6)|25)
	})

	t.Run("lines=3", func(t *testing.T) {
		run := func(maxHeight int, wantAscent, wantDescent fixed.Int26_6) {
			t.Run(fmt.Sprintf("height=%d", maxHeight), func(t *testing.T) {
				check(t, 3, maxHeight, wantAscent, wantDescent)
			})
		}
		run(20, fixed.I(5)|29, fixed.I(1)|14)
		run(30, fixed.I(8)|11, fixed.I(1)|53)
		run(40, fixed.I(10)|57, fixed.I(2)|28)
		run(50, fixed.I(13)|39, fixed.I(3)|2)
		run(60, fixed.I(16)|22, fixed.I(3)|42)
		run(70, fixed.I(19)|5, fixed.I(4)|17)
	})
}
