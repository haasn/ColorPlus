package colorplus

import "math"

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
func FromTemperature(T float64) Yxy {
    var x float64

    if (4000 <= T && T <= 7000) {
        x = -4.6070E9 / math.Pow(T, 3) + 2.9678E6 / math.Pow(T, 2) + 9.911E1 / T + 0.244063
    } else if (7000 < T && T <= 25000) {
        x = -2.0064E9 / math.Pow(T, 3) + 1.9018E6 / math.Pow(T, 2) + 2.4748E2 / T + 0.237040
    } else {
        panic("[FromTemperature] Color temperature out of range!")
    }

    y := -3 * x * x + 2.87 * x - 0.275

    return Yxy{1, x, y}
}
