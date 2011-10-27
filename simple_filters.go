package colorplus

// Default filters
var Identity = FilterSingle(func(in float64) float64 { // returns argument unmodified
    return in
})

var Invert = FilterSingle(func(in float64) float64 { // inverts an argument on the scale 0-1
    return 1 - in
})
