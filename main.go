package main

import (
	"os"
    "bytes"
    "bufio"
    "image/jpeg"
    "image/png"
    "flag"
    "fmt"
    "path"

    "github.com/icza/mjpeg"
)

var frames = flag.Int("f", 120, "number of `frames`")
var rate = flag.Int("r", 30, "frame `rate` (frames/second)")
var filename = flag.String("o", "simplex.mjpg", "output movie `filename`")

func main() {
    flag.Parse()

    width := 32
    height := 16

    aw, err := mjpeg.New(*filename, int32(width), int32(height), int32(*rate))
    if err != nil {
        fatal("cannot create movie", err)
    }
    defer aw.Close()

    f := NewFrame(width, height)
    for i := 0; i < *frames; i++ {
        m := f.Next()

        if true && i == 0 {
            p, err := os.OpenFile(fmt.Sprintf("frame%02d.png", i), os.O_WRONLY|os.O_CREATE, 0600)
            defer p.Close()

            err = png.Encode(p, m)
            if err != nil {
                fatal("cannot encode image", err)
            }
        }

        var b bytes.Buffer
        w := bufio.NewWriter(&b)
        err = jpeg.Encode(w, m, nil)
        if err != nil {
            fatal("cannot encode image", err)
        }

        err = aw.AddFrame(b.Bytes())
    }
}

func fatal(reason string, err error) {
    fmt.Fprintf(os.Stderr, "%s: %s: %v\n", path.Base(os.Args[0]), reason, err)
    os.Exit(1)
}
