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
