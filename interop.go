package colorplus

import (
    "image"
    "math"
)

// Our RGB type can implement image.Color - note that colors must already be pulled up as appropriate
func (in RGB) RGBA() (r, g, b, a uint32) {
    return uint32(math.Floor(in.R)), uint32(math.Floor(in.G)), uint32(math.Floor(in.B)), 0xFFFF
}

func ApplyToImage(i *image.RGBA, f FilterTripleProvider) {
    w, h := i.Rect.Max.X, i.Rect.Max.Y
    filter := Pipeline(16, true, Chain(f, Clamp{0, 1})).GetTriple()

    for x := 0; x < w; x++ {
        for y := 0; y < h; y++ {
            r, g, b, _ := i.At(x, y).RGBA()
            i.Set(x, y, filter(RGB{float64(r), float64(g), float64(b)}).(RGB))
        }
    }
}

func Pipeline(depth uint, full bool, filter FilterTripleProvider) FilterTripleProvider {
    return Chain(Pulldown{depth, full}, filter, Pullup{depth, full})
}
