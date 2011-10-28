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

func FuzzyCompareMatrix3x3(x, y matrix3x3, allow float64) bool {
    err := matrix3x3{math.Abs(x.M11 - y.M11), math.Abs(x.M12 - y.M12), math.Abs(x.M13 - y.M13),
                     math.Abs(x.M21 - y.M21), math.Abs(x.M22 - y.M22), math.Abs(x.M23 - y.M23),
                     math.Abs(x.M31 - y.M31), math.Abs(x.M32 - y.M32), math.Abs(x.M33 - y.M33)}
    
    return (err.M11 < allow) && (err.M12 < allow) && (err.M13 < allow) &&
           (err.M21 < allow) && (err.M22 < allow) && (err.M23 < allow) &&
           (err.M31 < allow) && (err.M32 < allow) && (err.M33 < allow)
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
    nSRGBEnc = namedFilter{SRGBCurve.GetEncoder(), "SRGBCurve.GetEncoder()"}
    nSRGBDec = namedFilter{SRGBCurve.GetDecoder(), "SRGBCurve.GetDecoder()"}
    nLStarEnc = namedFilter{LStarActual.GetEncoder(), "LStarCurve{LStarActual}.GetEncoder()"}
    nLStarDec = namedFilter{LStarActual.GetDecoder(), "LStarCurve{LStarActual}.GetDecoder()"}
    nGrayscale = namedFilter{Grayscale, "Grayscale"}
)
