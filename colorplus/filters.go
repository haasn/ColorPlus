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
type filterTripleChain []FilterTripleProvider
type filterSingleChain []FilterSingleProvider

// Function to concatenate filters together
func Chain(list ...FilterTripleProvider) (FilterTripleProvider) {
    return filterTripleChain(list)
}

func ChainSingle(list ...FilterSingleProvider) (FilterSingleProvider) {
    return filterSingleChain(list)
}

// Filter provider implementation for filter chains
func (ftc filterTripleChain) GetTriple() FilterTriple {
    cache := make([]FilterTriple, len(ftc))

    for i,f := range ftc {
        cache[i] = f.GetTriple()
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
