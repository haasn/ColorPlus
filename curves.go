package colorplus

import "math"

// A curve provider is used for gamma encoding/decoding
type CurveProvider interface {
    GetEncoder() FilterSingle
    GetDecoder() FilterSingle
}

// Some built-in curve types
type purePowerCurve struct {
    gamma float64
}

type sRGBCurve struct {}

type lStarCurve struct {
    mode LStarMode
}

// Generator functions / variables
func PurePowerCurve(gamma float64) CurveProvider {
    return purePowerCurve{gamma}
}

var (
    SRGBCurve CurveProvider = sRGBCurve{}
    LStarActual CurveProvider = lStarCurve{lStarActual}
    LStarIntent CurveProvider = lStarCurve{lStarIntent}
)

// L* mode determines whether to follow the intent of the CIE Standard, or its written value
type LStarMode bool

const (
    lStarActual LStarMode = true
    lStarIntent LStarMode = false
)

// Implementations
func (ppc purePowerCurve) GetEncoder() FilterSingle {
    gp := 1.0 / ppc.gamma
    return func(in float64) float64 {
        if (in > 0) {
            return math.Pow(in, gp)
        }
        return 0
    }
}

func (ppc purePowerCurve) GetDecoder() FilterSingle {
    return func(in float64) float64 {
        if (in > 0) {
            return math.Pow(in, ppc.gamma)
        }
        return 0
    }
}

// sRGB curve is a two-part curve with a linear segment for blacks
func (_ sRGBCurve) GetEncoder() FilterSingle {
    return func(in float64) float64 {
        if (in > 0.0031308) {
            return math.Pow(in, 1 / 2.4) * 1.055 - 0.055
        }
        return in * 12.92
    }
}

func (_ sRGBCurve) GetDecoder() FilterSingle {
    return func(in float64) float64 {
        if (in > 0.04045) {
            return math.Pow((in + 0.055) / 1.055, 2.4)
        }
        return in / 12.92
    }
}

// L* Curve is the curve used for most broadcasting systems
func (ls lStarCurve) GetEncoder() FilterSingle {
    E, K := ls.params()

    return func(in float64) float64 {
        if (in > E) {
            return math.Pow(in, 1.0 / 3.0) * 1.16 - 0.16
        }
        return in * K / 100.0
    }
}

func (ls lStarCurve) GetDecoder() FilterSingle {
    E, K := ls.params()
    crossover := E * K / 100.0

    return func(in float64) float64 {
        if (in > crossover) {
            return math.Pow((in + 0.16) / 1.16, 3)
        }
        return in * 100.0 / K
    }
}

func (ls lStarCurve) params() (E, K float64) {
    if (ls.mode == lStarActual) {
        E, K = 0.008856, 903.3
    } else {
        E, K = 216.0 / 24389.0, 24389.0 / 27.0
    }

    return
}
