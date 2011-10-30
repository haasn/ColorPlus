package colorplus

import "testing"

func TestLut(t *testing.T) {
    lut := Make3DLUT([]int32{8, 8, 8}, 8, RangeFull, RangeFull, EncodingBGR, EncodingBGR, SpacesRGB, SpacesRGB)

    lut.Assign(Chain(Grayscale, Invert), true)
    rgb := lut.GetOutputRaw(0).(RGB)

    FuzzyAssertTriple(RGB{255, 0, 0}, rgb, RGB{200, 200, 200}, allow, "lut.Assign(Chain(Gra yscale, Invert), true)", t)
}
