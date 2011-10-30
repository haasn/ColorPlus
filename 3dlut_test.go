package colorplus

import "testing"

func TestLut(t *testing.T) {
    lut := Make3DLUT([]int32{4, 4, 4}, 8, RangeFull, RangeFull, EncodingBGR, EncodingBGR, SpacesRGB, SpacesRGB)

    lut.Assign(Chain(Grayscale, Invert), true)
    rgb := lut.GetOutputRaw(lut.Offset(0, 0, 15)).(RGB)

    FuzzyAssertTriple(RGB{15, 0, 0}, rgb, RGB{200, 200, 200}, allow, "lut.Assign(Chain(Grayscale, Invert), true)", t)
}
