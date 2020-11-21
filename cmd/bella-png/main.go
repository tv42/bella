package main

import (
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"eagain.net/go/bella/internal/render"
	"golang.org/x/sys/unix"
)

func main() {
	prog := filepath.Base(os.Args[0])
	log.SetFlags(0)
	log.SetPrefix(prog + ": ")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", prog)
		fmt.Fprintf(flag.CommandLine.Output(), "  %s [OPTS] LINES.. >DEST_PNG\n", prog)
		fmt.Fprintln(flag.CommandLine.Output())
		flag.PrintDefaults()
	}

	var height int
	flag.IntVar(&height, "height", 64, "Height of image")
	flag.Parse()

	if flag.NArg() == 0 {
		log.Print("missing text to print")
		os.Exit(2)
	}

	if _, err := unix.IoctlGetTermios(unix.Stdout, unix.TCGETS); err == nil {
		// is a terminal
		log.Fatal("stdout is a terminal, refusing to output binary")
	}

	text := strings.Join(flag.Args(), "\n")
	img, err := render.Text(text,
		render.WithMaxHeight(height),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(os.Stdout, img); err != nil {
		log.Fatalf("png encode: %v", err)
	}
}
