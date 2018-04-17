// MIT License
//
// Copyright 2018 Jeremy Hall
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Command line Simplex Noise movie generator.
//
// When I use the VLC media player to play this movie on the mac I get a
// "Broken or missing index" error, but the movie plays OK if I press the "Play
// as is" button.

package main

import (
    "bufio"
    "bytes"
    "flag"
    "fmt"
    "image/jpeg"
    "image/png"
	"os"
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

        if false {
            // Save frame as png
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
        if err != nil {
            fatal("cannot add frame to movie", err)
        }
    }
}

func fatal(reason string, err error) {
    fmt.Fprintf(os.Stderr, "%s: %s: %v\n", path.Base(os.Args[0]), reason, err)
    os.Exit(1)
}
