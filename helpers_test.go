package colorplus

import (
    "math"
    "testing"
)

// fuzzy comparison for color triples
func FuzzyCompareTriple(x, y Triple, allow float64) bool {
	xa, xb, xc := x.Get()
	ya, yb, yc := y.Get()

	erra, errb, errc := math.Abs(xa-ya), math.Abs(xb-yb), math.Abs(xc-yc)

	return (erra < allow) && (errb < allow) && (errc < allow)
}

func FuzzyCompareSingle(x, y float64, allow float64) bool {
    return math.Abs(x - y) < allow
}

func FuzzyAssertTriple(in interface{}, res, want Triple, allow float64, name string, t *testing.T) {
    if (!FuzzyCompareTriple(res, want, allow)) {
        t.Errorf("%s(%v) = %v, want %v.", name, in, res, want)
    }
}

func FuzzyAssertSingle(in interface{}, res, want float64, allow float64, name string, t *testing.T) {
    if (!FuzzyCompareSingle(res, want, allow)) {
        t.Errorf("%s(%v) = %v, want %v.", name, in, res, want)
    }
}

const allow = 0.00001

// interface to encapsulate named filters
type namedFilter struct {
    filter interface{}
    name string
}

var (
    nIdentity = namedFilter{Identity, "Identity"}
    nInvert = namedFilter{Invert, "Invert"}
    nXYZtoYxy = namedFilter{XYZtoYxy, "XYZtoYxy"}
    nYxytoXYZ = namedFilter{YxytoXYZ, "YxytoXYZ"}
    nSRGBEnc = namedFilter{SRGBCurve{}.GetEncoder(), "SRGBCurve.GetEncoder()"}
    nSRGBDec = namedFilter{SRGBCurve{}.GetDecoder(), "SRGBCurve.GetDecoder()"}
    nLStarEnc = namedFilter{LStarCurve{LStarActual}.GetEncoder(), "LStarCurve{LStarActual}.GetEncoder()"}
    nLStarDec = namedFilter{LStarCurve{LStarActual}.GetDecoder(), "LStarCurve{LStarActual}.GetDecoder()"}
    nGrayscale = namedFilter{Grayscale, "Grayscale"}
)
