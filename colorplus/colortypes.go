package colorplus

// Usage notes: For the sake of compatibility, all of these structures are assumed to be normalized into the range 0-1
// The functions which convert to/from “dumb” color types (like image.ColorRGBA) will perform pullup as appropriate
// Similarly, the functions to Normalize/Denormalize XYZ values will do this as well
// If a struct is marked with a ', it is to be considered non-linear

// Some generic interfaces for dealing with n-channeled colors
type Triple interface {
    Get() (a, b, c float64)
    Make(a, b, c float64) Triple // The make function serves the following purpose: If the underlying structure type is not known,
                                 // one can use the Make() function to create a new value of that type without losing any parameters
}

// type Single float64

// CIE 1931 XYZ
type XYZ struct {
    X, Y, Z float64
}

func (in XYZ) Get() (a, b, c float64) {
    return in.X, in.Y, in.Z
}

func (_ XYZ) Make(a, b, c float64) Triple {
    return XYZ{a, b, c}
}

// CIE Yxy (derived from XYZ)
type Yxy struct {
    Y, x, y float64
}

func (in Yxy) Get() (a, b, c float64) {
    return in.Y, in.x, in.y
}

func (_ Yxy) Make(a, b, c float64) Triple {
    return Yxy{a, b, c}
}
