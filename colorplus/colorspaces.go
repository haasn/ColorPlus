package colorplus

// The general interface - this shouldn't be abused too much because the functions may do heavy conversions in the background
type Space interface {
    RedXYZ() XYZ
    GreenXYZ() XYZ
    BlueXYZ() XYZ
    WhiteXYZ() XYZ // White().Y doubles as reference luminosity
}

// The individual implementations
type SpaceXYZ struct {
    Red, Green, Blue, White XYZ
}

func (cs SpaceXYZ) RedXYZ() XYZ { return cs.Red }
func (cs SpaceXYZ) GreenXYZ() XYZ { return cs.Green }
func (cs SpaceXYZ) BlueXYZ() XYZ { return cs.Blue }
func (cs SpaceXYZ) WhiteXYZ() XYZ { return cs.White }

// Yxy
type SpaceYxy struct {
    Red, Green, Blue, White Yxy
}

func (cs SpaceYxy) RedXYZ() XYZ { return cs.Red.ToXYZ() }
func (cs SpaceYxy) GreenXYZ() XYZ { return cs.Green.ToXYZ() }
func (cs SpaceYxy) BlueXYZ() XYZ { return cs.Blue.ToXYZ() }
func (cs SpaceYxy) WhiteXYZ() XYZ { return cs.White.ToXYZ() }
