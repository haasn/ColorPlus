package colorplus

import "math"

// A curve provider is used for gamma encoding/decoding
type CurveProvider interface {
    GetEncoder() FilterSingle
    GetDecoder() FilterSingle
}

// Some built-in curve types
type PurePowerCurve struct {
    Gamma float64
}

type sRGBCurve byte

// L* mode determines whether to follow the intent of the CIE Standard, or its written value, represented here as bool
type LStarCurve bool

const (
    LStarActual LStarCurve = true
    LStarIntent LStarCurve = false
    SRGBCurve sRGBCurve = 0
)

// Implementations
func (ppc PurePowerCurve) GetEncoder() FilterSingle {
    gp := 1.0 / ppc.Gamma
    return func(in float64) float64 {
        if (in > 0) {
            return math.Pow(in, gp)
        }
        return 0
    }
}

func (ppc PurePowerCurve) GetDecoder() FilterSingle {
    return func(in float64) float64 {
        if (in > 0) {
            return math.Pow(in, ppc.Gamma)
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
func (ls LStarCurve) GetEncoder() FilterSingle {
    E, K := ls.params()

    return func(in float64) float64 {
        if (in > E) {
            return math.Pow(in, 1.0 / 3.0) * 1.16 - 0.16
        }
        return in * K / 100.0
    }
}

func (ls LStarCurve) GetDecoder() FilterSingle {
    E, K := ls.params()
    crossover := E * K / 100.0

    return func(in float64) float64 {
        if (in > crossover) {
            return math.Pow((in + 0.16) / 1.16, 3)
        }
        return in * 100.0 / K
    }
}

func (ls LStarCurve) params() (E, K float64) {
    if bool(ls) {
        E, K = 0.008856, 903.3
    } else {
        E, K = 216.0 / 24389.0, 24389.0 / 27.0
    }

    return
}
