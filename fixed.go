package main

import(
    "fmt"
    "math"
)

type Fixed float64

var min Fixed
var max Fixed

func record(f Fixed) {
    if f > max {
        max = f
    } else if f < min {
        min = f
    }
}

func IntToFixed(i int) Fixed {
    result := Fixed(i)
    record(result)
    return result
}

func (x Fixed) ToInt() int {
    return int(math.Round(float64(x)))
}

func FloatToFixed(f float64) Fixed {
    result := Fixed(f)
    record(result)
    return result
}

func (x Fixed) ToFloat() float64 {
    return float64(x)
}

func (x Fixed) Mul(y Fixed) Fixed {
    result := Fixed(x*y)
    record(result)
    return result
}

func (x Fixed) Div(y Fixed) Fixed {
    result := Fixed(x/y)
    record(result)
    return result
}

func (x Fixed) Abs() Fixed {
    result := x
    if x < 0 {
        result = -x
    }
    record(result)
    return result
}

func (x Fixed) Floor() Fixed {
    result := Fixed(math.Floor(x.ToFloat()))
    record(result)
    return result
}

func (x Fixed) Ceil() Fixed {
    result := Fixed(math.Ceil(x.ToFloat()))
    record(result)
    return result
}

func (x Fixed) Round() Fixed {
    result := Fixed(math.Round(x.ToFloat()))
    record(result)
    return result
}

func ReportLimits() {
    fmt.Printf("Fixed.min: %f\n", min.ToFloat())
    fmt.Printf("Fixed.max: %f\n", max.ToFloat())
}
