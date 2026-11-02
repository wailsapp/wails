// Package typeconstraints computes the structural restrictions used by the
// binding generator. Method and comparable requirements remain go/types' job.
package typeconstraints

import (
	"errors"
	"go/types"
)

var (
	ErrEmpty   = errors.New("empty structural type set")
	ErrInvalid = errors.New("invalid or cyclic type constraint")
	ErrComplex = errors.New("type constraint exceeds normalization limits")
)

const maxTerms = 100

type termSet struct {
	all   bool
	terms []*types.Term
}

// Terms returns a minimal union of exact and underlying-type terms. A nil slice
// with no error means unrestricted, while ErrEmpty means no structural matches.
// The result deliberately excludes method and comparable filtering.
func Terms(typ types.Type) ([]*types.Term, error) {
	w := walker{memo: make(map[types.Type]termSet), active: make(map[types.Type]bool)}
	set, err := w.expand(typ)
	if err != nil {
		return nil, err
	}
	if !set.all && len(set.terms) == 0 {
		return nil, ErrEmpty
	}
	return set.terms, nil
}

type walker struct {
	memo   map[types.Type]termSet
	active map[types.Type]bool
}

func (w *walker) expand(typ types.Type) (termSet, error) {
	if typ == nil || typ == types.Typ[types.Invalid] || w.active[typ] {
		return termSet{}, ErrInvalid
	}
	if result, ok := w.memo[typ]; ok {
		return result, nil
	}
	w.active[typ] = true
	defer delete(w.active, typ)
	result, err := w.resolve(types.Unalias(typ))
	if err == nil {
		w.memo[typ] = result
	}
	return result, err
}

func (w *walker) resolve(typ types.Type) (termSet, error) {
	if parameter, ok := typ.(*types.TypeParam); ok {
		return w.expand(parameter.Constraint())
	}
	switch underlying := typ.Underlying().(type) {
	case *types.Interface:
		result := termSet{all: true}
		for i := 0; i < underlying.NumEmbeddeds(); i++ {
			embedded, err := w.expand(underlying.EmbeddedType(i))
			if err != nil {
				return termSet{}, err
			}
			result = intersect(result, embedded)
		}
		return result, nil
	case *types.Union:
		return w.union(underlying)
	default:
		return termSet{terms: []*types.Term{types.NewTerm(false, typ)}}, nil
	}
}

func (w *walker) union(union *types.Union) (termSet, error) {
	var result termSet
	for i := 0; i < union.Len(); i++ {
		part, err := w.expandTerm(union.Term(i))
		if err != nil {
			return termSet{}, err
		}
		if part.all {
			result = part
		} else if !result.all {
			for _, candidate := range part.terms {
				result.terms = include(result.terms, candidate)
			}
		}
		if len(result.terms) > maxTerms {
			return termSet{}, ErrComplex
		}
	}
	return result, nil
}

func (w *walker) expandTerm(term *types.Term) (termSet, error) {
	if term.Type() == nil || term.Type() == types.Typ[types.Invalid] {
		return termSet{}, ErrInvalid
	}
	if !term.Tilde() {
		return w.expand(term.Type())
	}
	return termSet{terms: []*types.Term{term}}, nil
}

// covers reports whether every type in narrow is also in broad.
func covers(broad, narrow *types.Term) bool {
	if !broad.Tilde() {
		return !narrow.Tilde() && types.Identical(broad.Type(), narrow.Type())
	}
	return types.Identical(broad.Type().Underlying(), narrow.Type().Underlying())
}

func include(terms []*types.Term, candidate *types.Term) []*types.Term {
	for _, existing := range terms {
		if covers(existing, candidate) {
			return terms
		}
	}
	kept := terms[:0]
	for _, existing := range terms {
		if !covers(candidate, existing) {
			kept = append(kept, existing)
		}
	}
	return append(kept, candidate)
}

func intersect(a, b termSet) termSet {
	if a.all {
		return b
	}
	if b.all {
		return a
	}
	var result termSet
	for _, left := range a.terms {
		for _, right := range b.terms {
			if covers(left, right) {
				result.terms = include(result.terms, right)
			} else if covers(right, left) {
				result.terms = include(result.terms, left)
			}
		}
	}
	return result
}
