package colorplus

import "math"

// Default filters
var Identity = FilterSingle(func(in float64) float64 { // returns argument unmodified
    return in
})

var Invert = FilterSingle(func(in float64) float64 { // inverts an argument on the scale 0-1
    return 1 - in
})

// Some simple filter providers
type Clamp struct {
    Lower, Upper float64
}

func (cf Clamp) GetSingle() FilterSingle {
    return func(in float64) float64 {
        return math.Max(cf.Lower, math.Min(cf.Upper, in))
    }
}

// Modify the value range
type Scale struct {
    Lower, Upper float64
}

func (s Scale) GetSingle() FilterSingle {
    return func(in float64) float64 {
        return in * (s.Upper - s.Lower) + s.Lower
    }
}

// Swap channels around
type Swap []SwapMode
type SwapMode byte

const (
    AB SwapMode = 0
    AC SwapMode = 1
    BC SwapMode = 2
)

func (s Swap) GetTriple() FilterTriple {
    return func(in Triple) Triple {
        a, b, c := in.Get()
        for _,v := range s {
            switch (v) {
                case AB: a, b = b, a
                case AC: a, c = c, a
                case BC: b, c = c, b
                default: panic("Invalid swap mode!")
            }
        }
        return in.Make(a, b, c)
    }
}
