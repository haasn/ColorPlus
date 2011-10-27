package colorplus

import "math"

// fuzzy comparison for color triples
func FuzzyCompare(x, y Triple, allow float64) bool {
	xa, xb, xc := x.Get()
	ya, yb, yc := y.Get()

	erra, errb, errc := math.Abs(xa-ya), math.Abs(xb-yb), math.Abs(xc-yc)

	return (erra < allow) && (errb < allow) && (errc < allow)
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
)
