package colorplus

// Determine the luminance (brightness) of a color triple
func Luminance(in Triple) float64 {
    switch x := in.(type) {
        case XYZ: return x.Y
        case Yxy: return x.Y
        case RGB: return 0.2126 * x.R + 0.0722 * x.B + 0.7152 * x.G
        default: panic("[Luminance] Unsupported color type!")
    }
    return 0
}

// Calculate a white point from a color temperature
func FromTemperature(T float64) XYZ {
    var x float64

    if (4000 <= T && T <= 7000) {
        x = -4.6070E9 / (T*T*T) + 2.9678E6 / (T*T) + 9.911E1 / T + 0.244063
    } else if (7000 < T && T <= 25000) {
        x = -2.0064E9 / (T*T*T) + 1.9018E6 / (T*T) + 2.4748E2 / T + 0.237040
    } else {
        panic("[FromTemperature] Color temperature out of range!")
    }

    y := -3 * x * x + 2.87 * x - 0.275

    return Yxy{1, x, y}.ToXYZ()
}

// Chromatic adaptation
type ChromaticAdapter struct {
    Source, Destination XYZ
    Mode ScalingMode
}

type ScalingMode matrix3x3

var (
    Linear = ScalingMode{1, 0, 0, 0, 1, 0, 0, 0, 1}
    Bradford = ScalingMode{0.8951000, 0.2664000, -0.1614000, -0.7502000, 1.7135000, 0.0367000, 0.0389000, -0.0685000, 1.0296000} // recommended
    VonKries = ScalingMode{0.4002400, 0.7076000, -0.0808100, -0.2263000, 1.1653200, 0.0457000, 0.0000000, 0.0000000, 0.9182200}
)

func (ca ChromaticAdapter) GetTriple() FilterTriple {
    forw := matrix3x3(ca.Mode)
    back := forw.Inverse()

    cs := forw.Mul1x3(matrix1x3{ca.Source.X, ca.Source.Y, ca.Source.Z})
    cd := forw.Mul1x3(matrix1x3{ca.Destination.X, ca.Destination.Y, ca.Destination.Z})

    transform := back.Mul3x3(matrix3x3{cd.M1 / cs.M1, 0, 0, 0, cd.M2 / cs.M2, 0, 0, 0, cd.M3 / cs.M3}).Mul3x3(forw)

    return func(in Triple) Triple {
        var x XYZ
        switch v := in.(type) {
            case XYZ: x = v
            case Yxy: x = v.ToXYZ()
            default: panic("[ChromaticAdapter] Unsupported color type!")
        }
        res := transform.Mul1x3(matrix1x3{x.X, x.Y, x.Z})

        return XYZ{res.M1, res.M2, res.M3}
    }
}

// XYZ encoding/decoding
type XYZSpace Space // Regular Space will auto-apply the gamma transfer function as well

func (c XYZSpace) GetEncoder() FilterTriple {
    M := matrixFromColorSpace(Space(c)).Inverse()

    return func(in Triple) Triple {
        var x XYZ
        switch v := in.(type) {
            case XYZ: x = v
            case Yxy: x = v.ToXYZ()
            default: panic("[XYZSpace.GetEncoder] Unsupported color type!")
        }

        res := M.Mul1x3(matrix1x3{x.X, x.Y, x.Z})

        return RGB{res.M1, res.M2, res.M3}
    }
}

func (c XYZSpace) GetDecoder() FilterTriple {
    M := matrixFromColorSpace(Space(c))

    return func(in Triple) Triple {
        switch v := in.(type) {
            case RGB:
                res := M.Mul1x3(matrix1x3{v.R, v.G, v.B})
                return XYZ{res.M1, res.M2, res.M3}
            default: panic("[XYZSpace.GetDecoder] Unsupported color type!")
        }

        return nil
    }
}

// Space encoding/decoding, with gamma
func (c Space) GetEncoder() FilterTriple {
    if (c.Gamma == nil) {
        return XYZSpace(c).GetEncoder()
    } else {
        return Chain(XYZSpace(c).GetEncoder(), c.Gamma.GetEncoder()).GetTriple()
    }
    return nil // Why does the compiler need a return here?
}

func (c Space) GetDecoder() FilterTriple {
    if (c.Gamma == nil) {
        return XYZSpace(c).GetDecoder()
    } else {
        return Chain(c.Gamma.GetDecoder(), XYZSpace(c).GetDecoder()).GetTriple()
    }
    return nil
}

// Custom matrix logic
type matrix1x3 struct {
    M1, M2, M3 float64
}

type matrix3x3 struct {
    M11, M12, M13,
    M21, M22, M23,
    M31, M32 , M33 float64
}

func (m matrix3x3) Inverse() matrix3x3 {
    t1 := m.M33 * m.M22 - m.M32 * m.M23
    t2 := m.M33 * m.M12 - m.M32 * m.M13
    t3 := m.M23 * m.M12 - m.M22 * m.M13
    det := m.M11 * t1 - m.M21 * t2 + m.M31 * t3;

    return matrix3x3{t1, -t2, t3,
           -(m.M33 * m.M21 - m.M31 * m.M23), m.M33 * m.M11 - m.M31 * m.M13, -(m.M23 * m.M11 - m.M21 * m.M13),
           m.M32 * m.M21 - m.M31 * m.M22, -(m.M32 * m.M11 - m.M31 * m.M12), m.M22 * m.M11 - m.M21 * m.M12}.MulC(1 / det)
}

// M * c
func (m matrix3x3) MulC(c float64) matrix3x3 {
    return matrix3x3{m.M11 * c, m.M12 * c, m.M13 * c, m.M21 * c, m.M22 * c, m.M23 * c, m.M31 * c, m.M32 * c, m.M33 * c}
}

// M * M
func (a matrix3x3) Mul3x3(b matrix3x3) matrix3x3 {
    return matrix3x3{
        a.M11 * b.M11 + a.M12 * b.M21 + a.M13 * b.M31, // C11
        a.M11 * b.M12 + a.M12 * b.M22 + a.M13 * b.M32, // C12
        a.M11 * b.M13 + a.M12 * b.M23 + a.M13 * b.M33, // C13

        a.M21 * b.M11 + a.M22 * b.M21 + a.M23 * b.M31, // C21
        a.M21 * b.M12 + a.M22 * b.M22 + a.M23 * b.M32, // C22
        a.M21 * b.M13 + a.M22 * b.M23 + a.M23 * b.M33, // C23

        a.M31 * b.M11 + a.M32 * b.M21 + a.M33 * b.M31, // C21
        a.M31 * b.M12 + a.M32 * b.M22 + a.M33 * b.M32, // C22
        a.M31 * b.M13 + a.M32 * b.M23 + a.M33 * b.M33} // C23
}

func (a matrix3x3) Mul1x3(b matrix1x3) matrix1x3 {
    return matrix1x3{
        a.M11 * b.M1 + a.M12 * b.M2 + a.M13 * b.M3,
        a.M21 * b.M1 + a.M22 * b.M2 + a.M23 * b.M3,
        a.M31 * b.M1 + a.M32 * b.M2 + a.M33 * b.M3}
}

func matrixFromColorSpace(c Space) matrix3x3 {
    S := matrix3x3{c.Red.X, c.Green.X, c.Blue.X,
                   c.Red.Y, c.Green.Y, c.Blue.Y,
                   c.Red.Z, c.Green.Z, c.Blue.Z}.Inverse().Mul1x3(matrix1x3{c.White.X, c.White.Y, c.White.Z})
    
    return matrix3x3{S.M1 * c.Red.X, S.M2 * c.Green.X, S.M3 * c.Blue.X,
                     S.M1 * c.Red.Y, S.M2 * c.Green.Y, S.M3 * c.Blue.Y,
                     S.M1 * c.Red.Z, S.M2 * c.Green.Z, S.M3 * c.Blue.Z}
}
