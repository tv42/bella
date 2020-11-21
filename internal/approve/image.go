package approve

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func colorEq(a, b color.Color) bool {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
}

func Image(t testing.TB, img image.Image) error {
	if img == nil {
		return errors.New("image is nil")
	}
	callers := make([]uintptr, 1)
	n := runtime.Callers(2, callers)
	if n < 1 {
		return errors.New("unknown caller")
	}
	frames := runtime.CallersFrames(callers)
	frame, _ := frames.Next()

	if frame.File == "" {
		return fmt.Errorf("caller file not known: %+v", frame)
	}
	sourceDir := filepath.Dir(frame.File)
	testdata := path.Join(sourceDir, "testdata")
	testname := t.Name()
	// sanitize subtest names for filesystem use
	testname = strings.ReplaceAll(testname, "%", "%25")
	testname = strings.ReplaceAll(testname, ".", "%2E")

	newName := path.Join(testdata, testname+".new.png")
	if err := os.Mkdir(filepath.Dir(newName), 0755); err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("cannot create directory: %w", err)
	}
	newF, err := os.Create(newName)
	if err != nil {
		return fmt.Errorf("cannot open PNG for saving: %v", err)
	}
	defer newF.Close()
	if err := png.Encode(newF, img); err != nil {
		return fmt.Errorf("cannot save PNG: %v: %v", newName, err)
	}
	if err := newF.Close(); err != nil {
		return fmt.Errorf("cannot finish saving PNG: %v", err)
	}

	goodName := path.Join(testdata, testname+".good.png")
	goodF, err := os.Open(goodName)
	if err != nil {
		return fmt.Errorf("cannot open good file: %v", err)
	}
	defer goodF.Close()
	good, err := png.Decode(goodF)
	if err != nil {
		return fmt.Errorf("cannot load good PNG: %v: %v", goodName, err)
	}

	if g, e := good.Bounds().Size(), img.Bounds().Size(); !g.Eq(e) {
		return fmt.Errorf("size mismatch: %v != %v", g, e)
	}

	if err := strictEq(img, good); err == nil {
		return nil
	}
	// try a fuzzy match
	//
	// If you know of a Go library implementing a perceptual diff (not
	// just a perceptual hash), please tell me!
	cmd := exec.Command("perceptualdiff", newName, goodName)
	if buf, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("perceptualdiff failed:\n%s%s", buf, err)
	}
	return nil
}

// strictEq is a strict pixel by pixel comparison. It assumes the
// images are the same size.
func strictEq(img, good image.Image) error {
	off := img.Bounds().Min.Sub(good.Bounds().Min)
	for y := good.Bounds().Min.Y; y < good.Bounds().Max.Y; y++ {
		for x := good.Bounds().Min.X; x < good.Bounds().Max.X; x++ {
			colorGood := good.At(x, y)
			p := image.Point{X: x, Y: y}.Add(off)
			colorGot := img.At(p.X, p.Y)
			if !colorEq(colorGood, colorGot) {
				return fmt.Errorf("pixel difference at %v: %v != %v", p, colorGot, colorGood)
			}
		}
	}
	return nil
}
