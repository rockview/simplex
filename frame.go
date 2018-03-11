package main

import(
    "image"
    "image/color"
    "image/draw"
    "math"
)

type Frame struct {
    width, height int
    mode int
    xy_scale float64
    z_step float64
    z_depth float64
    z_state float64
    hue_options float64
    hue_state float64
    min, max float64
}

func NewFrame(width, height int) *Frame {
    return &Frame{
        width: width, height: height,
        mode: 2,
        xy_scale:  8.0/64.0,
        z_step: 0.0125,
        z_depth: 512,
        z_state: 0,
        hue_options: 0.005,
        hue_state: 0,
        min: 0.0001, max: 0.0001,
    }
}

func (f *Frame) Next() image.Image {
    m := image.NewRGBA(image.Rect(0, 0, f.width, f.height))

    for y := 0; y < f.height; y++ {
        sy := float64(y)*f.xy_scale;
        for x := 0; x < f.width; x++ {
            sx := float64(x)*f.xy_scale;
            n1 := noise(sx, sy, f.z_state)
            n2 := noise(sx, sy, f.z_state - f.z_depth)
            n := ((f.z_depth - f.z_state)*n1 + f.z_state*n2)/f.z_depth
            if n > f.max {
                f.max = n
            }
            if n < f.min {
                f.min = n
            }
            n += math.Abs(f.min)
            n /= (f.max + math.Abs(f.min))
            r := image.Rect(x, y, x + 1, y + 1)
            switch f.mode {
            case 1:
                hue := int((f.hue_options + n)*90.0 + 0.5)%90
                draw.Draw(m, r, &image.Uniform{getColor(hue)}, image.ZP, draw.Src)
            case 2:
                hue := int((f.hue_state + n)*90.0 + 0.5)%90
                draw.Draw(m, r, &image.Uniform{getColor(hue)}, image.ZP, draw.Src)
            case 3:
                hue := int(f.hue_state*90.0 + 0.5)%90
                draw.Draw(m, r, &image.Uniform{getColor(hue)}, image.ZP, draw.Src)
            }
        }
    }

    f.z_state = math.Mod(f.z_state + f.z_step, f.z_depth)
    f.hue_state = math.Mod(f.hue_state + f.hue_options, 1.0)

    return m
}

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

func getColor(hue int) color.Color {
    return color.RGBA{hues[hue][0], hues[hue][1], hues[hue][2], 0xff}
}
