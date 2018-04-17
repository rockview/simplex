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

// Simplex Noise frame generator.
//
// Inspired by code I found here: https://github.com/bikerglen/beagle/tree/master/projects/led-panel-v01/software

package main

import(
    "fmt"
    "image"
    "image/color"
    "image/draw"

    "github.com/rockview/fixed"
)

// Frame state
type Frame struct {
    width       int
    height      int
    mode        int
    xy_scale    fixed.Fixed
    z_step      fixed.Fixed
    z_depth     fixed.Fixed
    z_state     fixed.Fixed
    hue_options fixed.Fixed
    hue_state   fixed.Fixed
    min         fixed.Fixed
    max         fixed.Fixed
}

// NewFrame creates a new Frame instance
func NewFrame(width, height int) *Frame {
    return &Frame{
        width:          width,
        height:         height,
        mode:           2,
        xy_scale:       fixed.FromFloat(8.0/64.0),
        z_step:         fixed.FromFloat(0.0125),
        z_depth:        fixed.FromInt(512),
        z_state:        fixed.Zero,
        hue_options:    fixed.FromFloat(0.005),
        hue_state:      fixed.Zero,
        min:            fixed.FromFloat(0.0001),
        max:            fixed.FromFloat(0.0001),
    }
}

// Next draws the next frame
func (f *Frame) Next() image.Image {
    m := image.NewRGBA(image.Rect(0, 0, f.width, f.height))
    fixed90 := fixed.FromInt(90)

    for y := 0; y < f.height; y++ {
        sy := fixed.Mul(fixed.FromInt(y), f.xy_scale)
        for x := 0; x < f.width; x++ {
            // Draw pixel with computed hue
            sx := fixed.Mul(fixed.FromInt(x), f.xy_scale)
            n1 := noise(sx, sy, f.z_state)
            n2 := noise(sx, sy, f.z_state - f.z_depth)
            m1 := fixed.Mul(f.z_depth - f.z_state, n1)
            m2 := fixed.Mul(f.z_state, n2)
            n := fixed.Div(m1 + m2, f.z_depth)
            if n > f.max {
                f.max = n
            }
            if n < f.min {
                f.min = n
            }
            n += fixed.Abs(f.min)
            n = fixed.Div(n, f.max + fixed.Abs(f.min))
            r := image.Rect(x, y, x + 1, y + 1)
            var hue int
            switch f.mode {
            case 1:
                hue = fixed.ToInt(fixed.Round(fixed.Mul(f.hue_options + n, fixed90)))%90
            case 2:
                hue = fixed.ToInt(fixed.Round(fixed.Mul(f.hue_state + n, fixed90)))%90
            case 3:
                hue = fixed.ToInt(fixed.Round(fixed.Mul(f.hue_state, fixed90)))%90
            default:
                panic(fmt.Sprintf("invalid mode: %d", f.mode))
            }
            draw.Draw(m, r, &image.Uniform{getColor(hue)}, image.ZP, draw.Src)
        }
    }

    f.z_state = fixed.Mod(f.z_state + f.z_step, f.z_depth)
    f.hue_state = fixed.Mod(f.hue_state + f.hue_options, fixed.One)

    //f.dump()

    return m
}

// RGB hue color values
var hues = [...][3]uint8{
	{0xff, 0x00, 0x00},	//  0
	{0xff, 0x00, 0x11},	//  1
	{0xff, 0x00, 0x22},	//  2
	{0xff, 0x00, 0x33},	//  3
	{0xff, 0x00, 0x44},	//  4
	{0xff, 0x00, 0x55},	//  5
	{0xff, 0x00, 0x66},	//  6
	{0xff, 0x00, 0x77},	//  7
	{0xff, 0x00, 0x88},	//  8
	{0xff, 0x00, 0x99},	//  9
	{0xff, 0x00, 0xaa},	// 10
	{0xff, 0x00, 0xbb},	// 11
	{0xff, 0x00, 0xcc},	// 12
	{0xff, 0x00, 0xdd},	// 13
	{0xff, 0x00, 0xee},	// 14
	{0xff, 0x00, 0xff},	// 15
	{0xee, 0x00, 0xff},	// 16
	{0xdd, 0x00, 0xff},	// 17
	{0xcc, 0x00, 0xff},	// 18
	{0xbb, 0x00, 0xff},	// 19
	{0xaa, 0x00, 0xff},	// 20
	{0x99, 0x00, 0xff},	// 21
	{0x88, 0x00, 0xff},	// 22
	{0x77, 0x00, 0xff},	// 23
	{0x66, 0x00, 0xff},	// 24
	{0x55, 0x00, 0xff},	// 25
	{0x44, 0x00, 0xff},	// 26
	{0x33, 0x00, 0xff},	// 27
	{0x22, 0x00, 0xff},	// 28
	{0x11, 0x00, 0xff},	// 29
	{0x00, 0x00, 0xff},	// 30
	{0x00, 0x11, 0xff},	// 31
	{0x00, 0x22, 0xff},	// 32
	{0x00, 0x33, 0xff},	// 33
	{0x00, 0x44, 0xff},	// 34
	{0x00, 0x55, 0xff},	// 35
	{0x00, 0x66, 0xff},	// 36
	{0x00, 0x77, 0xff},	// 37
	{0x00, 0x88, 0xff},	// 38
	{0x00, 0x99, 0xff},	// 39
	{0x00, 0xaa, 0xff},	// 40
	{0x00, 0xbb, 0xff},	// 41
	{0x00, 0xcc, 0xff},	// 42
	{0x00, 0xdd, 0xff},	// 43
	{0x00, 0xee, 0xff},	// 44
	{0x00, 0xff, 0xff},	// 45
	{0x00, 0xff, 0xee},	// 46
	{0x00, 0xff, 0xdd},	// 47
	{0x00, 0xff, 0xcc},	// 48
	{0x00, 0xff, 0xbb},	// 49
	{0x00, 0xff, 0xaa},	// 50
	{0x00, 0xff, 0x99},	// 51
	{0x00, 0xff, 0x88},	// 52
	{0x00, 0xff, 0x77},	// 53
	{0x00, 0xff, 0x66},	// 54
	{0x00, 0xff, 0x55},	// 55
	{0x00, 0xff, 0x44},	// 56
	{0x00, 0xff, 0x33},	// 57
	{0x00, 0xff, 0x22},	// 58
	{0x00, 0xff, 0x11},	// 59
	{0x00, 0xff, 0x00},	// 60
	{0x11, 0xff, 0x00},	// 61
	{0x22, 0xff, 0x00},	// 62
	{0x33, 0xff, 0x00},	// 63
	{0x44, 0xff, 0x00},	// 64
	{0x55, 0xff, 0x00},	// 65
	{0x66, 0xff, 0x00},	// 66
	{0x77, 0xff, 0x00},	// 67
	{0x88, 0xff, 0x00},	// 68
	{0x99, 0xff, 0x00},	// 69
	{0xaa, 0xff, 0x00},	// 70
	{0xbb, 0xff, 0x00},	// 71
	{0xcc, 0xff, 0x00},	// 72
	{0xdd, 0xff, 0x00},	// 73
	{0xee, 0xff, 0x00},	// 74
	{0xff, 0xff, 0x00},	// 75
	{0xff, 0xee, 0x00},	// 76
	{0xff, 0xdd, 0x00},	// 77
	{0xff, 0xcc, 0x00},	// 78
	{0xff, 0xbb, 0x00},	// 79
	{0xff, 0xaa, 0x00},	// 80
	{0xff, 0x99, 0x00},	// 81
	{0xff, 0x88, 0x00},	// 82
	{0xff, 0x77, 0x00},	// 83
	{0xff, 0x66, 0x00},	// 84
	{0xff, 0x55, 0x00},	// 85
	{0xff, 0x44, 0x00},	// 86
	{0xff, 0x33, 0x00},	// 87
	{0xff, 0x22, 0x00},	// 88
	{0xff, 0x11, 0x00},	// 89
}

// getColor returns a color with the given hue
func getColor(hue int) color.Color {
    return color.RGBA{hues[hue][0], hues[hue][1], hues[hue][2], 0xff}
}

// dump dumps the frame state for debugging
func (f *Frame) dump() {
    fmt.Printf("---\n")
    fmt.Printf("width:       %d\n", f.width)
    fmt.Printf("height:      %d\n", f.height)
    fmt.Printf("mode:        %d\n", f.mode)
    fmt.Printf("xy_scale:    %f\n", fixed.ToFloat(f.xy_scale))
    fmt.Printf("z_step:      %f\n", fixed.ToFloat(f.z_step))
    fmt.Printf("z_depth:     %f\n", fixed.ToFloat(f.z_depth))
    fmt.Printf("z_state:     %f\n", fixed.ToFloat(f.z_state))
    fmt.Printf("hue_options: %f\n", fixed.ToFloat(f.hue_options))
    fmt.Printf("hue_state:   %f\n", fixed.ToFloat(f.hue_state))
    fmt.Printf("min:         %f\n", fixed.ToFloat(f.min))
    fmt.Printf("max:         %f\n", fixed.ToFloat(f.max))
}
