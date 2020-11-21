package render

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Option option

type config struct {
	maxHeight int
	imageFn   ImageFunc
	fontFn    FontFunc
}

type option func(*config)

// WithMaxHeight sets maximum height for rendered text, in pixels.
//
// This option is mandatory.
func WithMaxHeight(h int) Option {
	if h <= 0 {
		panic(fmt.Sprintf("WithMaxHeight: impossible height %d", h))
	}
	opt := func(cfg *config) {
		cfg.maxHeight = h
	}
	return opt
}

type ImageFunc func(image.Rectangle) draw.Image

// WithImageFunc sets the function used to create an image of the
// desired size.
//
// By default, image.NewGray is used.
func WithImageFunc(fn ImageFunc) Option {
	opt := func(cfg *config) {
		cfg.imageFn = fn
	}
	return opt
}

// FontFunc returns a new font face of the given point size.
//
// Points are 1/72 of an inch.
// https://en.wikipedia.org/wiki/Point_(typography)
type FontFunc func(size float64) (font.Face, error)

// WithFont sets the font used.
//
// By default, the Go Regular font is used.
func WithFont(fn FontFunc) Option {
	opt := func(cfg *config) {
		cfg.fontFn = fn
	}
	return opt
}

// WithFontOpenType is a convenience helper for WithFont that uses
// an OpenType font.
//
// Unlike opentype.NewFace, a FaceOptions value must be provided.
// FaceOptions.Size is overwritten by this function.
func WithFontOpenType(f *opentype.Font, opts *opentype.FaceOptions) Option {
	fn := func(size float64) (font.Face, error) {
		opts.Size = size
		face, err := opentype.NewFace(f, opts)
		if err != nil {
			return nil, err
		}
		return face, nil
	}
	return WithFont(fn)
}

// Text renders text as large as possible, without exceeding the
// maximum size.
func Text(text string, opts ...Option) (image.Image, error) {
	var cfg config
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.maxHeight == 0 {
		panic("render.Text: missing option WithMaxHeight")
	}
	if cfg.imageFn == nil {
		cfg.imageFn = func(rect image.Rectangle) draw.Image {
			return image.NewGray(rect)
		}
	}
	if cfg.fontFn == nil {
		ttf, err := opentype.Parse(goregular.TTF)
		if err != nil {
			return nil, fmt.Errorf("error parsing font Go Regular: %w", err)
		}
		opt := WithFontOpenType(ttf, &opentype.FaceOptions{
			DPI:     72,
			Hinting: font.HintingFull,
		})
		opt(&cfg)
	}
	return renderText(text, &cfg)
}

func renderText(text string, cfg *config) (image.Image, error) {
	if strings.Index(text, "\n") >= 0 {
		return nil, errors.New("multi-line text not supported yet") // TODO
	}
	face, err := binsearch(cfg.maxHeight, cfg.fontFn)
	if err != nil {
		return nil, err
	}
	defer face.Close()

	d := &font.Drawer{
		// Dst will be set later, once we know the right bounds

		Src:  image.NewUniform(color.Black),
		Face: face,
		// Adjust starting point to avoid negative bounding box
		// coordinates.
		Dot: fixed.Point26_6{X: 0, Y: face.Metrics().Ascent},
	}
	bounds, _ := d.BoundString(text)
	if bounds.Max.Y.Ceil() > cfg.maxHeight {
		return nil, fmt.Errorf("bounds spilled: %v..%v", bounds.Min, bounds.Max)
	}
	d.Dst = cfg.imageFn(image.Rect(0, 0, bounds.Max.X.Ceil(), cfg.maxHeight))

	// it starts out as all black, set a better background
	draw.Draw(d.Dst, d.Dst.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	d.DrawString(text)
	return d.Dst, nil
}
