package colorplus

// Filters are just functions which transform colors
type FilterTriple func(Triple) Triple
type FilterSingle func(float64) float64

type FilterTripleProvider interface {
    GetTriple() FilterTriple
}

type FilterSingleProvider interface {
    GetSingle() FilterSingle
}

// Two filters chained together
type filterTripleChain []interface{} // can be triple or single
type filterSingleChain []FilterSingleProvider
type filterMultiplex struct {
    a, b, c FilterSingleProvider
}

// Function to concatenate filters together
func Chain(list ...interface{}) (FilterTripleProvider) {
    return filterTripleChain(list)
}

func ChainSingle(list ...FilterSingleProvider) (FilterSingleProvider) {
    return filterSingleChain(list)
}

// Multiplex function to apply single filters to each channel of a triple
func Multiplex(a, b, c FilterSingleProvider) (FilterTripleProvider) {
    return filterMultiplex{a, b, c}
}

// Filter provider implementation for filter chains
func (ftc filterTripleChain) GetTriple() FilterTriple {
    cache := make([]FilterTriple, len(ftc))

    for i,f := range ftc {
        switch x := f.(type) {
            case FilterTripleProvider: cache[i] = x.GetTriple()
            case FilterSingleProvider: cache[i] = x.GetSingle().GetTriple()
        }
    }

    return func(in Triple) Triple {
        for _,v := range cache {
            in = v(in)
        }
        return in
    }
}

func (fsc filterSingleChain) GetSingle() FilterSingle {
    cache := make([]FilterSingle, len(fsc))

    for i,f := range fsc {
        cache[i] = f.GetSingle()
    }

    return func(in float64) float64 {
        for _,v := range cache {
            in = v(in)
        }
        return in
    }
}

// Filter provider implementation for multiplex
func (fm filterMultiplex) GetTriple() FilterTriple {
    fa, fb, fc := fm.a.GetSingle(), fm.b.GetSingle(), fm.c.GetSingle()

    return func(in Triple) Triple {
        a, b, c := in.Get()
        return in.Make(fa(a), fb(b), fc(c))
    }
}

// Filter provider implementation for basic filters
func (f FilterTriple) GetTriple() FilterTriple {
    return f
}

func (f FilterSingle) GetTriple() FilterTriple {
    return func(in Triple) Triple {
        a, b, c := in.Get()
        return in.Make(f(a), f(b), f(c))
    }
}

func (f FilterSingle) GetSingle() FilterSingle {
    return f
}
