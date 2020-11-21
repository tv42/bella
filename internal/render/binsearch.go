package render

import (
	"log"
	"math"
	"sort"

	"golang.org/x/image/font"
)

func binsearch(maxHeight int, fontFn FontFunc) (font.Face, error) {
	const debug = false
	// One inch tall text, 2.54cm, should definitely be more than your
	// label maker tape is wide.
	const max = 72.0
	// Adjustment size, since this is floats not ints. Having this be
	// tiny adds very little cost to the logarithmic search, and helps
	// with my perfectionism.
	const epsilon = 0.01

	// sort.Search looks for smallest i=[0,n) for which fn(i) is true.
	// We want to find the largest size which still fit inside
	// constraints. Search for smallest point size that fails to fit,
	// then go one down from that.
	n := int(math.Ceil(max / epsilon))
	iToPoints := func(i int) float64 {
		// constrain input to >=0 so we don't have to worry about
		// negative numbers in the degenerate case where
		// sort.Search()==0.
		return math.Max(float64(i), 0) * epsilon
	}
	var faceErr error
	matchFn := func(i int) bool {
		size := iToPoints(i)
		face, err := fontFn(size)
		if err != nil {
			// awkward early abort
			faceErr = err
			return true
		}
		metrics := face.Metrics()
		// As per
		// https://developer.apple.com/library/archive/documentation/TextFonts/Conceptual/CocoaTextArchitecture/Art/glyph_metrics_2x.png
		// (via https://godoc.org/golang.org/x/image/font#Metrics),
		height := (metrics.Ascent + metrics.Descent).Ceil()
		face.Close()

		if height > maxHeight {
			if debug {
				log.Printf("size=%-5g big %v > %v", size, height, maxHeight)
			}
			return true
		}
		if debug {
			log.Printf("size=%-5g ok  %v <= %v", size, height, maxHeight)
		}
		return false
	}
	tooBig := sort.Search(n, matchFn)
	if faceErr != nil {
		return nil, faceErr
	}
	size := iToPoints(tooBig - 1)
	face, err := fontFn(size)
	if err != nil {
		return nil, err
	}
	return face, nil
}
