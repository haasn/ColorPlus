package colorplus

import "testing"

type testPair struct {
    nfilter namedFilter
    input, output Triple
}

var tests = []testPair{
    // Invert()
    testPair{nInvert, XYZ{0, 0, 0}, XYZ{1, 1, 1}},
    testPair{nInvert, XYZ{0.5, 0.5, 0.5}, XYZ{0.5, 0.5, 0.5}},
    testPair{nInvert, XYZ{0.2, 0.9, 0.4}, XYZ{0.8, 0.1, 0.6}},
    testPair{nInvert, Yxy{0.3, 0.3, 0.6}, Yxy{0.7, 0.7, 0.4}},

    // Identity()
    testPair{nIdentity, XYZ{0.2, 0.5, 0.9}, XYZ{0.2, 0.5, 0.9}},

    // Conversion between XYZ and Yxy
    testPair{nXYZtoYxy, XYZ{0.2, 0.4, 0.8}, Yxy{0.4, 0.14286, 0.28571}},
    testPair{nXYZtoYxy, XYZ{1, 1, 1}, Yxy{1, 1.0 / 3, 1.0 / 3}},
    testPair{nYxytoXYZ, Yxy{1, 0.3127, 0.3290}, XYZ{0.950455, 1, 1.089057}},
}

func TestFilters(t *testing.T) {
    for _, tp := range tests {
        res := tp.nfilter.filter.GetTriple()(tp.input)

        if (!FuzzyCompare(res, tp.output, allow)) {
            t.Errorf("%s(%v) = %v, want %v.", tp.nfilter.name, tp.input, res, tp.output)
        }
    }
}