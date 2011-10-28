package colorplus

import "math"

type Space struct {
    Red, Green, Blue, White XYZ // White().Y doubles as reference luminosity
    Gamma CurveProvider // optional and may be nil
}

// Functions for generating color spaces from xy values
func SpaceFromxy(rx, ry, gx, gy, bx, by float64, White Yxy, g CurveProvider) Space {
    return Space{Yxy{1, rx, ry}.ToXYZ(), Yxy{1, gx, gy}.ToXYZ(), Yxy{1, bx, by}.ToXYZ(), White.ToXYZ(), g}
}

func SpaceFromExisting(space Space, g CurveProvider) Space {
    return Space{space.Red, space.Green, space.Blue, space.White, g}
}

// Area calculation, useful for comparing sizes
func (s Space) Area() float64 {
    r, g, b := s.Red.ToYxy(), s.Green.ToYxy(), s.Blue.ToYxy()

    x := math.Sqrt(math.Pow(g.x - r.x, 2) + math.Pow(g.y - r.y, 2))
    y := math.Sqrt(math.Pow(g.x - b.x, 2) + math.Pow(g.y - b.y, 2))
    z := math.Sqrt(math.Pow(r.x - b.x, 2) + math.Pow(r.y - b.y, 2))

    return math.Sqrt((x + y - z) * (x - y + z) * (y + z - x) * (x + y + z)) / 4
}

// Default white points
var (
    PointA = Yxy{1, 0.44757, 0.40735}
    PointB = Yxy{1, 0.34842, 0.35161}
    PointC = Yxy{1, 0.31006, 0.31616}
    PointD50 = Yxy{1, 0.34567, 0.35850}
    PointD55 = Yxy{1, 0.33242, 0.34743}
    PointD65 = Yxy{1, 0.31272661468101209, 0.32902313032606195}
    PointD75 = Yxy{1, 0.29902, 0.31485}
    PointE = Yxy{1, 1.0 / 3.0, 1.0 / 3.0}
    PointF1 = Yxy{1, 0.31310, 0.33727}
    PointF2 = Yxy{1, 0.372080, 0.37529}
    PointF3 = Yxy{1, 0.40910, 0.39430}
    PointF4 = Yxy{1, 0.44018, 0.40329}
    PointF5 = Yxy{1, 0.31379, 0.34531}
    PointF6 = Yxy{1, 0.37790, 0.38835}
    PointF7 = Yxy{1, 0.31292, 0.32933}
    PointF8 = Yxy{1, 0.34588, 0.35875}
    PointF9 = Yxy{1, 0.37417, 0.37281}
    PointF10 = Yxy{1, 0.34609, 0.35986}
    PointF11 = Yxy{1, 0.38052, 0.37713}
    PointF12 = Yxy{1, 0.43695, 0.40441}
    PointZero = Yxy{1, 0, 0}
)

// Default color spaces
var (
    SpaceZero = SpaceFromxy(0, 0, 0, 0, 0, 0, PointZero, nil)
    SpaceBT709 = SpaceFromxy(0.64, 0.33, 0.30, 0.60, 0.15, 0.06, PointD65, PurePowerCurve{2.35})
    SpaceROMM = SpaceFromxy(0.7347, 0.2653, 0.1596, 0.8404, 0.0366, 0.0001, PointD50, nil)
    SpaceAdobeRGB98 = SpaceFromxy(0.64, 0.33, 0.21, 0.71, 0.15, 0.06, PointD65, nil)
    SpaceAppleRGB = SpaceFromxy(0.625, 0.34, 0.28, 0.595, 0.115, 0.07, PointD65, PurePowerCurve{1.8})
    SpaceNTSC_53 = SpaceFromxy(0.67, 0.33, 0.21, 0.71, 0.14, 0.08, PointC, nil)
    SpaceNTSC_87 = SpaceFromxy(0.63, 0.34, 0.31, 0.595, 0.155, 0.07, PointD65, PurePowerCurve{2.2})
    SpaceSECAM = SpaceFromxy(0.64, 0.33, 0.29, 0.60, 0.15, 0.06, PointD65, PurePowerCurve{2.8})
    SpaceAdobeWideRGB = SpaceFromxy(0.735, 0.265, 0.115, 0.826, 0.157, 0.018, PointD50, nil)
    SpaceCIE1931 = SpaceFromxy(0.7347, 0.2653, 0.2738, 0.7174, 0.1666, 0.0089, PointE, nil)
    SpaceACES = SpaceFromxy(0.73470, 0.26530, 0, 1, 0.00010, -0.07700, Yxy{1, 0.32168, 0.33767}, nil)

    // some aliases
    SpaceNone = SpaceZero

    SpaceHDTV = SpaceFromExisting(SpaceBT709, LStarCurve{LStarActual})
    SpacesRGB = SpaceFromExisting(SpaceBT709, SRGBCurve{})
    SpacescRGB = SpacesRGB
    SpacemadVR = SpaceFromExisting(SpaceBT709, PurePowerCurve{2.2})

    SpaceProPhotoRGB = SpaceROMM

    SpaceFCC1953 = SpaceNTSC_53
    SpaceBT470_M = SpaceNTSC_53

    SpaceSMPTE_C = SpaceNTSC_87
    SpaceSMPTE_RP_145 = SpaceNTSC_87
    SpaceSMPTE_170M = SpaceNTSC_87

    SpacePAL = SpaceSECAM
    SpaceEBU_Tech_3213 = SpaceSECAM
    SpaceBT470_B = SpaceSECAM
)
