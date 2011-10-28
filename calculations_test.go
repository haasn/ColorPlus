package colorplus

import "testing"

func TestCalculations(t *testing.T) {
    FuzzyAssertTriple(6503.6, FromTemperature(6503.6), PointD65, allow * 100, "FromTemperature", t) // not as precise, lots of conversions
    FuzzyAssertSingle("", SpacesRGB.Area(), 0.11205, allow, "SpacesRGB.Area", t)

    // Matrix testing
    m := matrixFromColorSpace(SpacesRGB)
    a := matrix3x3{0.4124564, 0.35757, 0.1804374, 0.2126728, 0.71515, 0.0721749, 0.0193338, 0.11919, 0.9503040}

    if !FuzzyCompareMatrix3x3(a, m, allow) {
        t.Errorf("matrixFromColorSpace(SpacesRGB) = %v, want %v.", m, a)
    }
}
