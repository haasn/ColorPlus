package colorplus

// From XYZ
func (in XYZ) ToYxy() Yxy {
    lum := in.X + in.Y + in.Z

    return Yxy{in.Y, in.X / lum, in.Y / lum}
}

// From Yxy
func (in Yxy) ToXYZ() XYZ {
    return XYZ{in.Y * in.x / in.y, in.Y, in.Y * (1 - in.x - in.y) / in.y}
}

// Conversion filters
var XYZtoYxy = FilterTriple(func(in Triple) Triple {
    return in.(XYZ).ToYxy()
})

var YxytoXYZ = FilterTriple(func(in Triple) Triple {
    return in.(Yxy).ToXYZ()
})
