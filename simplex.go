package main

import(
    "github.com/rockview/fixed"
)

type Grad struct {
    x, y, z fixed.Fixed
}

var grad3 = [...]Grad{
     Grad{ 1, 1, 0},
     Grad{-1, 1, 0},
     Grad{ 1,-1, 0},
     Grad{-1,-1, 0},
     Grad{ 1, 0, 1},
     Grad{-1, 0, 1},
     Grad{ 1, 0,-1},
     Grad{-1, 0,-1},
     Grad{ 0, 1, 1},
     Grad{ 0,-1, 1},
     Grad{ 0, 1,-1},
     Grad{ 0,-1,-1},
}

var perm [512]int
var permMod12 [512]int

// 3D simplex noise
func noise(xin, yin, zin fixed.Fixed) fixed.Fixed {
    // Skew the input space to determine which simplex cell we're in
    s := fixed.Mul(xin + yin + zin, fixed.Third)   // Very nice and simple skew factor for 3D
    i := fixed.Floor(xin + s)
    j := fixed.Floor(yin + s)
    k := fixed.Floor(zin + s)
    t := fixed.Mul(i + j + k, fixed.Sixth)
    X0 := i - t // Unskew the cell origin back to (x,y,z) space
    Y0 := j - t
    Z0 := k - t
    x0 := xin - X0  // The x,y,z distances from the cell origin
    y0 := yin - Y0
    z0 := zin - Z0

    // For the 3D case, the simplex shape is a slightly irregular tetrahedron.
    // Determine which simplex we are in.
    var i1, j1, k1 int  // Offsets for second corner of simplex in (i,j,k) coords
    var i2, j2, k2 int  // Offsets for third corner of simplex in (i,j,k) coords
    if (x0 >=  y0) {
        if (y0 >=  z0) {
            // X Y Z order
            i1 = 1; j1 = 0; k1 = 0
            i2 = 1; j2 = 1; k2 = 0
        } else if (x0 >=  z0) {
            // X Z Y order
            i1 = 1; j1 = 0; k1 = 0
            i2 = 1; j2 = 0; k2 = 1
        } else {
            // Z X Y order
            i1 = 0; j1 = 0; k1 = 1
            i2 = 1; j2 = 0; k2 = 1
        }
    } else {
        // x0 < y0
        if (y0 < z0) {
            // Z Y X order
            i1 = 0; j1 = 0; k1 = 1
            i2 = 0; j2 = 1; k2 = 1
        } else if(x0 < z0) {
            // Y Z X order
            i1 = 0; j1 = 1; k1 = 0
            i2 = 0; j2 = 1; k2 = 1
        } else {
            // Y X Z order
            i1 = 0; j1 = 1; k1 = 0
            i2 = 1; j2 = 1; k2 = 0
        }
    }

    // A step of (1,0,0) in (i,j,k) means a step of (1-c,-c,-c) in (x,y,z),
    // a step of (0,1,0) in (i,j,k) means a step of (-c,1-c,-c) in (x,y,z), and
    // a step of (0,0,1) in (i,j,k) means a step of (-c,-c,1-c) in (x,y,z), where
    // c = 1/6.
    x1 := x0 - fixed.FromInt(i1) + fixed.Sixth  // Offsets for second corner in (x,y,z) coords
    y1 := y0 - fixed.FromInt(j1) + fixed.Sixth
    z1 := z0 - fixed.FromInt(k1) + fixed.Sixth
    x2 := x0 - fixed.FromInt(i2) + fixed.Third  // Offsets for third corner in (x,y,z) coords
    y2 := y0 - fixed.FromInt(j2) + fixed.Third
    z2 := z0 - fixed.FromInt(k2) + fixed.Third
    x3 := x0 - fixed.One + fixed.Half // Offsets for last corner in (x,y,z) coords
    y3 := y0 - fixed.One + fixed.Half
    z3 := z0 - fixed.One + fixed.Half

    // Work out the hashed gradient indices of the four simplex corners
    ii := int(i)&255
    jj := int(j)&255
    kk := int(k)&255
    gi0 := permMod12[ii + perm[jj + perm[kk]]]
    gi1 := permMod12[ii + i1 + perm[jj + j1 + perm[kk + k1]]]
    gi2 := permMod12[ii + i2 + perm[jj + j2 + perm[kk + k2]]]
    gi3 := permMod12[ii + 1 + perm[jj + 1 + perm[kk + 1]]]

    var n0, n1, n2, n3 fixed.Fixed  // Noise contributions from the four corners
    t0 := fixed.Half - fixed.Mul(x0, x0) - fixed.Mul(y0, y0) - fixed.Mul(z0, z0)
    if (t0 >= 0) {
        t0 = fixed.Mul(t0, t0)
        n0 = fixed.Mul(fixed.Mul(t0, t0), dot(&grad3[gi0], x0, y0, z0))
    }
    t1 := fixed.Half - fixed.Mul(x1, x1) - fixed.Mul(y1, y1) - fixed.Mul(z1, z1)
    if (t1 >= 0) {
        t1 = fixed.Mul(t1, t1)
        n1 = fixed.Mul(fixed.Mul(t1, t1), dot(&grad3[gi1], x1, y1, z1))
    }
    t2 := fixed.Half - fixed.Mul(x2, x2) - fixed.Mul(y2, y2) - fixed.Mul(z2, z2)
    if (t2 >= 0) {
        t2 = fixed.Mul(t2, t2)
        n2 = fixed.Mul(fixed.Mul(t2, t2), dot(&grad3[gi2], x2, y2, z2))
    }
    t3 := fixed.Half - fixed.Mul(x3, x3) - fixed.Mul(y3, y3) - fixed.Mul(z3, z3)
    if (t3 >= 0) {
        t3 = fixed.Mul(t3, t3)
        n3 = fixed.Mul(fixed.Mul(t3, t3), dot(&grad3[gi3], x3, y3, z3))
    }

    // Add contributions from each corner to get the final noise value.
    // The result is scaled to stay just inside [-1,1]
    return fixed.Mul(fixed.FromFloat(32.0), (n0 + n1 + n2 + n3))
}

// Package intialization
func init() {
    p := [...]int{
        151, 160, 137, 91, 90, 15, 131, 13, 201, 95, 96, 53, 194, 233, 7, 225,
        140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23, 190, 6, 148,
        247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32,
        57, 177, 33, 88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175,
        74, 165, 71, 134, 139, 48, 27, 166, 77, 146, 158, 231, 83, 111, 229, 122,
        60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244, 102, 143, 54,
        65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169,
        200, 196, 135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64,
        52, 217, 226, 250, 124, 123, 5, 202, 38, 147, 118, 126, 255, 82, 85, 212,
        207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42, 223, 183, 170, 213,
        119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
        129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104,
        218, 246, 97, 228, 251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241,
        81, 51, 145, 235, 249, 14, 239, 107, 49, 192, 214, 31, 181, 199, 106, 157,
        184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254, 138, 236, 205, 93,
        222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180,
    }
    for i := 0; i < len(perm); i++ {
        perm[i] = p[i&255]
        permMod12[i] = perm[i]%12
    }
}

func dot(g *Grad, x, y, z fixed.Fixed) fixed.Fixed {
    return fixed.Mul(g.x, x) + fixed.Mul(g.y, y) + fixed.Mul(g.z, z)
}
