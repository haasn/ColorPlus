package colorplus

import "math"

// Returns argument unmodified
var Identity = FilterSingle(func(in float64) float64 {
    return in
})

// Inverts the input on the scale 0-1
var Invert = FilterSingle(func(in float64) float64 {
    return 1 - in
})

// Clamp the value to a given range
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
    AB SwapMode = iota
    AC
    BC
)

func (s Swap) GetTriple() FilterTriple {
    return func(in Triple) Triple {
        a, b, c := in.Get()
        for _,v := range s {
            switch (v) {
                case AB: a, b = b, a
                case AC: a, c = c, a
                case BC: b, c = c, b
                default: panic("[Swap] Invalid swap mode!")
            }
        }
        return in.Make(a, b, c)
    }
}

// Grayscale (Note: Only works on RGB because XYZ has unknown white point)
var Grayscale = FilterTriple(func(in Triple) Triple {
    switch x := in.(type) {
        case RGB:
            luma := Luminance(in)

            return in.Make(luma, luma, luma)
        default:
            panic("[Grayscale] Unsupported color type!")
    }
    return nil
})

// Range pullup
type Pullup struct {
    Depth uint
    FullRange bool // false = limited range
}

func (p Pullup) GetSingle() FilterSingle {
    bot, lim := calcLimits(p.Depth, p.FullRange)

    return func(in float64) float64 {
        return in * lim + bot
    }
}

// Range pulldown
type Pulldown Pullup

func (p Pulldown) GetSingle() FilterSingle {
    bot, lim := calcLimits(p.Depth, p.FullRange)

    return func(in float64) float64 {
        return (in - bot) / lim
    }
}
