package colorplus

import "testing"

func TestCalculations(t *testing.T) {
    FuzzyAssert(6504, FromTemperature(6504), PointD65, allow * 10, "FromTemperature", t)
}
