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

    // Multiplexing
    testPair{namedFilter{Multiplex(Invert, Identity, Invert), "Multiplex(Invert, Identity, Invert)"}, XYZ{0.2, 0.4, 0.4}, XYZ{0.8, 0.4, 0.6}},

    // Gamma transfer functions
    testPair{namedFilter{PurePowerCurve(2.2).GetDecoder(), "PurePowerCurve(2.2).GetDecoder()"},
        XYZ{0.2, 0, -1}, XYZ{0.028991, 0, 0}},
    
    testPair{namedFilter{PurePowerCurve(2.2).GetEncoder(), "PurePowerCurve(2.2).GetEncoder()"},
        XYZ{0.2, 0, -1}, XYZ{0.481156, 0, 0}},

    // sRGB
    testPair{nSRGBEnc, XYZ{0.001, 0.5, -1}, XYZ{0.01292, 0.735356, -12.92}},
    testPair{nSRGBDec, XYZ{0.01292, 0.735356, -12.92}, XYZ{0.001, 0.5, -1}},
}

func TestFilters(t *testing.T) {
    for _, tp := range tests {
        res := tp.nfilter.filter.GetTriple()(tp.input)

        if (!FuzzyCompare(res, tp.output, allow)) {
            t.Errorf("%s(%v) = %v, want %v.", tp.nfilter.name, tp.input, res, tp.output)
        }
    }
}
