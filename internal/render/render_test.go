package render_test

import (
	"image"
	"testing"

	"eagain.net/go/bella/internal/approve"
	"eagain.net/go/bella/internal/render"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

func mustPanic(t testing.TB, fn func()) {
	defer func() {
		p := recover()
		if p == nil {
			t.Fatal("expected a panic")
		}
		t.Logf("got expected panic: %v", p)
	}()
	fn()
}

func TestTextOptsNoHeight(t *testing.T) {
	mustPanic(t, func() { render.Text("foo") })
}

func TestText(t *testing.T) {
	ttf, err := opentype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("opentype.Parse: %v", err)
	}
	run := func(text string, maxHeight int) {
		t.Run(text, func(t *testing.T) {
			img, err := render.Text(
				text,
				render.WithMaxHeight(maxHeight),
				render.WithFontOpenType(
					ttf,
					&opentype.FaceOptions{
						DPI:     180,
						Hinting: font.HintingFull,
					},
				),
			)
			if err != nil {
				t.Fatalf("render.Text: %v", err)
			}
			if g, e := img.Bounds(), image.Rect(0, 0, img.Bounds().Max.X, maxHeight); g != e {
				t.Errorf("wrong bounds: %v != %v", g, e)
			}
			if err := approve.Image(t, img); err != nil {
				t.Errorf("not approved: %v", err)
			}
		})
	}

	run("Hello, world", 42)
	run(".", 64)
	run("tiny", 10)
	run("multi\nline", 64)
	run("way\ntoo\nsmall\nto\nread", 64)
}
