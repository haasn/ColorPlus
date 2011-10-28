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
    testPair{namedFilter{PurePowerCurve{2.2}.GetDecoder(), "PurePowerCurve{2.2}.GetDecoder()"},
        XYZ{0.2, 0, -1}, XYZ{0.028991, 0, 0}},
    
    testPair{namedFilter{PurePowerCurve{2.2}.GetEncoder(), "PurePowerCurve{2.2}.GetEncoder()"},
        XYZ{0.2, 0, -1}, XYZ{0.481156, 0, 0}},

    // sRGB
    testPair{nSRGBEnc, XYZ{0.001, 0.5, -1}, XYZ{0.01292, 0.735356, -12.92}},
    testPair{nSRGBDec, XYZ{0.01292, 0.735356, -12.92}, XYZ{0.001, 0.5, -1}},

    // Round trip 1
    testPair{namedFilter{Chain(Identity, Invert, Invert, nLStarEnc.filter, nLStarDec.filter),"RoundTrip1"},
        XYZ{0.2, 0.5, 0.8}, XYZ{0.2, 0.5, 0.8}},
    
    // Clamp and scale
    testPair{namedFilter{Clamp{0, 1}, "Clamp{0, 1}"}, XYZ{-2, 0.3, 1.5}, XYZ{0, 0.3, 1}},
    testPair{namedFilter{Scale{0.2, 0.8}, "Scale{0.2, 0.8}"}, XYZ{0, 1, 0.5}, XYZ{0.2, 0.8, 0.5}},

    // Swap
    testPair{namedFilter{Swap{AB, BC}, "Swap{AB, BC}"}, XYZ{0.2, 0.4, 0.8}, XYZ{0.4, 0.8, 0.2}},

    // Grayscale
    testPair{nGrayscale, RGB{0.2, 0.4, 0.8}, RGB{0.38636, 0.38636, 0.38636}},

    // XYZ Encoding/Decoding
    testPair{namedFilter{Chain(SpacesRGB.GetDecoder(), SpacesRGB.GetEncoder()), "Chain(SpacesRGB.GetDecoder(), SpacesRGB.GetEncoder())"},
        RGB{1.0, 1.0, 1.0}, RGB{1.0, 1.0, 1.0}},

    // White point check
    testPair{namedFilter{SpacesRGB.GetDecoder(), "SpacesRGB.GetDecoder()"}, RGB{1.0, 1.0, 1.0}, PointD65},

    // Chromatic adaptation
    testPair{namedFilter{ChromaticAdapter{PointD65, PointD50, Bradford}, "ChromaticAdapter{PointD65, PointD50, Bradford}"}, PointD65, PointD50},
}

func TestFilters(t *testing.T) {
    for _, tp := range tests {
        res := Chain(tp.nfilter.filter).GetTriple()(tp.input) // re-using Chain() for the type assertion logic
        FuzzyAssertTriple(tp.input, res, tp.output, allow, tp.nfilter.name, t)
    }
}
