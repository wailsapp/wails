package typeconstraints

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"slices"
	"strings"
	"testing"
)

const declarations = `package fixture
type Number int
type Alias = Number
type Exact interface { Number | string }
type Broad interface { ~int | ~string }
type Intersection interface { Exact; Broad }
type Absorbed interface { Number | Broad }
type Empty interface { ~int; ~string }
type Any interface {}
type Methods interface { M() }
type Comparable interface { comparable; ~int | ~[]byte }
type Generic[T interface { Intersection }] struct{}
type Repeat interface { Broad; Broad }
type Only interface { ~int; Number }
type Universal interface { Any | int }
type Nonempty interface { Empty | int }
type Reverse interface { Broad | Number }
`

func checkedPackage(t testing.TB) *types.Package {
	t.Helper()
	files := token.NewFileSet()
	file, err := parser.ParseFile(files, "fixture.go", declarations, 0)
	if err != nil {
		t.Fatal(err)
	}
	pkg, err := new(types.Config).Check("fixture", files, []*ast.File{file}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return pkg
}

func TestTerms(t *testing.T) {
	pkg := checkedPackage(t)
	for _, tt := range []struct {
		name string
		want []string
	}{
		{"Number", []string{"fixture.Number"}},
		{"Alias", []string{"fixture.Number"}},
		{"Exact", []string{"fixture.Number", "string"}},
		{"Broad", []string{"~int", "~string"}},
		{"Intersection", []string{"fixture.Number", "string"}},
		{"Absorbed", []string{"~int", "~string"}},
		{"Any", nil}, {"Methods", nil},
		// Structural terms deliberately leave comparable filtering to go/types.
		{"Comparable", []string{"~int", "~[]byte"}},
		{"Repeat", []string{"~int", "~string"}},
		{"Only", []string{"fixture.Number"}},
		{"Universal", nil},
		{"Nonempty", []string{"int"}},
		{"Reverse", []string{"~int", "~string"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			terms, err := Terms(pkg.Scope().Lookup(tt.name).Type())
			if err != nil {
				t.Fatal(err)
			}
			var got []string
			for _, term := range terms {
				got = append(got, term.String())
			}
			slices.Sort(got)
			slices.Sort(tt.want)
			if !slices.Equal(got, tt.want) {
				t.Fatalf("got %v; want %v", got, tt.want)
			}
		})
	}
	parameter := pkg.Scope().Lookup("Generic").Type().(*types.Named).TypeParams().At(0)
	if terms, err := Terms(parameter); err != nil || len(terms) != 2 {
		t.Fatalf("generic parameter: %v, %v", terms, err)
	}
	if _, err := Terms(pkg.Scope().Lookup("Empty").Type()); !errors.Is(err, ErrEmpty) {
		t.Fatalf("empty intersection: %v", err)
	}
}

func TestInvalidAndLimits(t *testing.T) {
	for _, typ := range []types.Type{nil, types.Typ[types.Invalid], types.NewUnion([]*types.Term{types.NewTerm(true, types.Typ[types.Invalid])})} {
		if _, err := Terms(typ); err == nil {
			t.Fatal("invalid constraint accepted")
		}
	}
	named := types.NewNamed(types.NewTypeName(token.NoPos, nil, "Cycle", nil), nil, nil)
	named.SetUnderlying(types.NewInterfaceType(nil, []types.Type{named}))
	if _, err := Terms(named); !errors.Is(err, ErrInvalid) {
		t.Fatalf("cycle: %v", err)
	}
	var terms []*types.Term
	for i := 0; i <= maxTerms; i++ {
		terms = append(terms, types.NewTerm(false, types.NewArray(types.Typ[types.Int], int64(i))))
	}
	if _, err := Terms(types.NewUnion(terms)); !errors.Is(err, ErrComplex) {
		t.Fatalf("term limit: %v", err)
	}
	var nested types.Type = types.Typ[types.Int]
	for i := 0; i < 1001; i++ {
		nested = types.NewInterfaceType(nil, []types.Type{nested})
	}
	if terms, err := Terms(nested); err != nil || len(terms) != 1 || !types.Identical(terms[0].Type(), types.Typ[types.Int]) {
		t.Fatalf("deep interface chain: %v, %v", terms, err)
	}
}

func generatedConstraint(data []byte) (*types.Interface, []types.Type) {
	basics := []types.Type{types.Typ[types.Int], types.Typ[types.String], types.Typ[types.Float64], types.NewSlice(types.Typ[types.Byte])}
	candidates := slices.Clone(basics)
	for _, basic := range basics {
		candidates = append(candidates, types.NewNamed(types.NewTypeName(token.NoPos, nil, "Named", nil), basic, nil))
	}
	var clauses []types.Type
	for _, mask := range data {
		var terms []*types.Term
		for i, basic := range basics {
			if mask&(1<<i) == 0 {
				continue
			}
			if mask&(1<<(i+4)) == 0 {
				terms = append(terms, types.NewTerm(false, candidates[i+4]))
			} else {
				terms = append(terms, types.NewTerm(true, basic))
			}
		}
		if len(terms) > 0 {
			clauses = append(clauses, types.NewUnion(terms))
		}
	}
	return types.NewInterfaceType(nil, clauses).Complete(), candidates
}

func FuzzStructuralMembership(f *testing.F) {
	f.Add([]byte{0xff, 0x15})
	f.Add([]byte{0x11, 0x22})
	f.Add([]byte{})
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 16 {
			return
		}
		constraint, candidates := generatedConstraint(data)
		terms, err := Terms(constraint)
		if err != nil && !errors.Is(err, ErrEmpty) {
			t.Fatal(err)
		}
		for _, candidate := range candidates {
			member := err == nil && len(terms) == 0
			for _, term := range terms {
				if term.Tilde() {
					member = member || types.Identical(term.Type(), candidate.Underlying())
				} else {
					member = member || types.Identical(term.Type(), candidate)
				}
			}
			if want := types.Satisfies(candidate, constraint); member != want {
				t.Fatalf("membership %s in %s: got %v, want %v", candidate, constraint, member, want)
			}
		}
	})
}

func BenchmarkTerms(b *testing.B) {
	pkg := checkedPackage(b)
	for _, name := range strings.Fields("Number Broad Intersection Absorbed Repeat Comparable") {
		typ := pkg.Scope().Lookup(name).Type()
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := Terms(typ); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
