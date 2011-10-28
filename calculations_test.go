package colorplus

import "testing"

func TestCalculations(t *testing.T) {
    FuzzyAssertTriple(6504, FromTemperature(6504), PointD65, allow * 10, "FromTemperature", t)
    FuzzyAssertSingle("", SpacesRGB.Area(), 0.11205, allow, "SpacesRGB.Area", t)
}
