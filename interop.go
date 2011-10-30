package colorplus

import (
    "image"
    "math"
)

// Our RGB type can implement image.Color - note that colors must already be pulled up as appropriate
func (in RGB) RGBA() (r, g, b, a uint32) {
    return uint32(in.R), uint32(in.G), uint32(in.B), 0xFFFF
}

func ApplyToImage(i *image.RGBA, f FilterTripleProvider) {
    w, h := i.Rect.Max.X, i.Rect.Max.Y
    filter := Pipeline(32, true, Chain(f, Clamp{0, 1})).GetTriple()

    for x := 0; x < w; x++ {
        for y := 0; y < h; y++ {
            r, g, b, _ := i.At(x, y).RGBA()
            post := filter(RGB{float64(r), float64(g), float64(b)}).(RGB)
            post.R, post.G, post.B = math.Floor(post.R), math.Floor(post.G), math.Floor(post.B)
            i.Set(x, y, post)
        }
    }
}

func Pipeline(depth uint, full bool, filter FilterTripleProvider) FilterTripleProvider {
    return Chain(Pulldown{depth, full}, filter, Pullup{depth, full})
}
