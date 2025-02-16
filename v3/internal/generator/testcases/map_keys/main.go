package main

import (
	_ "embed"
	"encoding"
	"encoding/json"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct{}

type NonTextMarshaler struct{}

type ValueTextMarshaler struct{}

func (ValueTextMarshaler) MarshalText() ([]byte, error) { return nil, nil }

type PointerTextMarshaler struct{}

func (*PointerTextMarshaler) MarshalText() ([]byte, error) { return nil, nil }

type JsonTextMarshaler struct{}

func (JsonTextMarshaler) MarshalJSON() ([]byte, error) { return nil, nil }
func (JsonTextMarshaler) MarshalText() ([]byte, error) { return nil, nil }

type CustomInterface interface {
	MarshalText() ([]byte, error)
}

type EmbeddedInterface interface {
	encoding.TextMarshaler
}

type EmbeddedInterfaces interface {
	json.Marshaler
	encoding.TextMarshaler
}

type BasicConstraint interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~string
}

type BadTildeConstraint interface {
	int | ~struct{} | string
}

type GoodTildeConstraint interface {
	int | ~struct{} | string
	MarshalText() ([]byte, error)
}

type NonBasicConstraint interface {
	ValueTextMarshaler | *PointerTextMarshaler
}

type PointableConstraint interface {
	ValueTextMarshaler | PointerTextMarshaler
}

type MixedConstraint interface {
	uint | ~string | ValueTextMarshaler | *PointerTextMarshaler
}

type InterfaceConstraint interface {
	comparable
	encoding.TextMarshaler
}

type PointerConstraint[T comparable] interface {
	*T
	encoding.TextMarshaler
}

type EmbeddedValue struct{ ValueTextMarshaler }
type EmbeddedValuePtr struct{ *ValueTextMarshaler }
type EmbeddedPointer struct{ PointerTextMarshaler }
type EmbeddedPointerPtr struct{ *PointerTextMarshaler }

type EmbeddedCustomInterface struct{ CustomInterface }
type EmbeddedOriginalInterface struct{ encoding.TextMarshaler }

type WrongType bool
type WrongAlias = bool
type StringType string
type StringAlias = string
type IntType int
type IntAlias = int

type ValueType ValueTextMarshaler
type ValuePtrType *ValueTextMarshaler
type ValueAlias = ValueTextMarshaler
type ValuePtrAlias = *ValueTextMarshaler

type PointerType PointerTextMarshaler
type PointerPtrType *PointerTextMarshaler
type PointerAlias = PointerTextMarshaler
type PointerPtrAlias = *PointerTextMarshaler

type InterfaceType encoding.TextMarshaler
type InterfacePtrType *encoding.TextMarshaler
type InterfaceAlias = encoding.TextMarshaler
type InterfacePtrAlias = *encoding.TextMarshaler

// type ComparableCstrAlias[R comparable] = R
// type ComparableCstrPtrAlias[R comparable] = *R
// type BasicCstrAlias[S BasicConstraint] = S
// type BasicCstrPtrAlias[S BasicConstraint] = *S
// type BadTildeCstrAlias[T BadTildeConstraint] = T
// type BadTildeCstrPtrAlias[T BadTildeConstraint] = *T
// type GoodTildeCstrAlias[U GoodTildeConstraint] = U
// type GoodTildeCstrPtrAlias[U GoodTildeConstraint] = *U
// type NonBasicCstrAlias[V NonBasicConstraint] = V
// type NonBasicCstrPtrAlias[V NonBasicConstraint] = *V
// type PointableCstrAlias[W PointableConstraint] = W
// type PointableCstrPtrAlias[W PointableConstraint] = *W
// type MixedCstrAlias[X MixedConstraint] = X
// type MixedCstrPtrAlias[X MixedConstraint] = *X
// type InterfaceCstrAlias[Y InterfaceConstraint] = Y
// type InterfaceCstrPtrAlias[Y InterfaceConstraint] = *Y
// type PointerCstrAlias[R comparable, Z PointerConstraint[R]] = Z
// type PointerCstrPtrAlias[R comparable, Z PointerConstraint[R]] = *Z

type Maps[R comparable, S BasicConstraint, T BadTildeConstraint, U GoodTildeConstraint, V NonBasicConstraint, W PointableConstraint, X MixedConstraint, Y InterfaceConstraint, Z PointerConstraint[R]] struct {
	Bool    map[bool]int      // Reject
	Int     map[int]int       // Accept
	Uint    map[uint]int      // Accept
	Float   map[float32]int   // Reject
	Complex map[complex64]int // Reject
	Byte    map[byte]int      // Accept
	Rune    map[rune]int      // Accept
	String  map[string]int    // Accept

	IntPtr     map[*int]int       // Reject
	UintPtr    map[*uint]int      // Reject
	FloatPtr   map[*float32]int   // Reject
	ComplexPtr map[*complex64]int // Reject
	StringPtr  map[*string]int    // Reject

	NTM    map[NonTextMarshaler]int      // Reject
	NTMPtr map[*NonTextMarshaler]int     // Reject
	VTM    map[ValueTextMarshaler]int    // Accept
	VTMPtr map[*ValueTextMarshaler]int   // Accept
	PTM    map[PointerTextMarshaler]int  // Reject
	PTMPtr map[*PointerTextMarshaler]int // Accept
	JTM    map[JsonTextMarshaler]int     // Accept, hide
	JTMPtr map[*JsonTextMarshaler]int    // Accept, hide

	A     map[any]int                     // Reject
	APtr  map[*any]int                    // Reject
	TM    map[encoding.TextMarshaler]int  // Accept, hide
	TMPtr map[*encoding.TextMarshaler]int // Reject
	CI    map[CustomInterface]int         // Accept, hide
	CIPtr map[*CustomInterface]int        // Reject
	EI    map[EmbeddedInterface]int       // Accept, hide
	EIPtr map[*EmbeddedInterface]int      // Reject

	EV     map[EmbeddedValue]int       // Accept
	EVPtr  map[*EmbeddedValue]int      // Accept
	EVP    map[EmbeddedValuePtr]int    // Accept
	EVPPtr map[*EmbeddedValuePtr]int   // Accept
	EP     map[EmbeddedPointer]int     // Reject
	EPPtr  map[*EmbeddedPointer]int    // Accept
	EPP    map[EmbeddedPointerPtr]int  // Accept
	EPPPtr map[*EmbeddedPointerPtr]int // Accept

	ECI    map[EmbeddedCustomInterface]int    // Accept
	ECIPtr map[*EmbeddedCustomInterface]int   // Accept
	EOI    map[EmbeddedOriginalInterface]int  // Accept
	EOIPtr map[*EmbeddedOriginalInterface]int // Accept

	WT   map[WrongType]int   // Reject
	WA   map[WrongAlias]int  // Reject
	ST   map[StringType]int  // Accept
	SA   map[StringAlias]int // Accept
	IntT map[IntType]int     // Accept
	IntA map[IntAlias]int    // Accept

	VT     map[ValueType]int      // Reject
	VTPtr  map[*ValueType]int     // Reject
	VPT    map[ValuePtrType]int   // Reject
	VPTPtr map[*ValuePtrType]int  // Reject
	VA     map[ValueAlias]int     // Accept
	VAPtr  map[*ValueAlias]int    // Accept
	VPA    map[ValuePtrAlias]int  // Accept, hide
	VPAPtr map[*ValuePtrAlias]int // Reject

	PT     map[PointerType]int      // Reject
	PTPtr  map[*PointerType]int     // Reject
	PPT    map[PointerPtrType]int   // Reject
	PPTPtr map[*PointerPtrType]int  // Reject
	PA     map[PointerAlias]int     // Reject
	PAPtr  map[*PointerAlias]int    // Accept
	PPA    map[PointerPtrAlias]int  // Accept, hide
	PPAPtr map[*PointerPtrAlias]int // Reject

	IT     map[InterfaceType]int      // Accept, hide
	ITPtr  map[*InterfaceType]int     // Reject
	IPT    map[InterfacePtrType]int   // Reject
	IPTPtr map[*InterfacePtrType]int  // Reject
	IA     map[InterfaceAlias]int     // Accept, hide
	IAPtr  map[*InterfaceAlias]int    // Reject
	IPA    map[InterfacePtrAlias]int  // Reject
	IPAPtr map[*InterfacePtrAlias]int // Reject

	TPR    map[R]int  // Soft reject
	TPRPtr map[*R]int // Soft reject
	TPS    map[S]int  // Accept, hide
	TPSPtr map[*S]int // Soft reject
	TPT    map[T]int  // Soft reject
	TPTPtr map[*T]int // Soft reject
	TPU    map[U]int  // Accept, hide
	TPUPtr map[*U]int // Soft reject
	TPV    map[V]int  // Accept, hide
	TPVPtr map[*V]int // Soft reject
	TPW    map[W]int  // Soft reject
	TPWPtr map[*W]int // Accept, hide
	TPX    map[X]int  // Accept, hide
	TPXPtr map[*X]int // Soft reject
	TPY    map[Y]int  // Accept, hide
	TPYPtr map[*Y]int // Soft reject
	TPZ    map[Z]int  // Accept, hide
	TPZPtr map[*Z]int // Soft reject

	// GAR    map[ComparableCstrAlias[R]]int    // Soft reject
	// GARPtr map[ComparableCstrPtrAlias[R]]int // Soft reject
	// GAS    map[BasicCstrAlias[S]]int         // Accept, hide
	// GASPtr map[BasicCstrPtrAlias[S]]int      // Soft reject
	// GAT    map[BadTildeCstrAlias[T]]int      // Soft reject
	// GATPtr map[BadTildeCstrPtrAlias[T]]int   // Soft reject
	// GAU    map[GoodTildeCstrAlias[U]]int     // Accept, hide
	// GAUPtr map[GoodTildeCstrPtrAlias[U]]int  // Soft reject
	// GAV    map[NonBasicCstrAlias[V]]int      // Accept, hide
	// GAVPtr map[NonBasicCstrPtrAlias[V]]int   // Soft reject
	// GAW    map[PointableCstrAlias[W]]int     // Soft reject
	// GAWPtr map[PointableCstrPtrAlias[W]]int  // Accept, hide
	// GAX    map[MixedCstrAlias[X]]int         // Accept, hide
	// GAXPtr map[MixedCstrPtrAlias[X]]int      // Soft reject
	// GAY    map[InterfaceCstrAlias[Y]]int     // Accept, hide
	// GAYPtr map[InterfaceCstrPtrAlias[Y]]int  // Soft reject
	// GAZ    map[PointerCstrAlias[R, Z]]int    // Accept, hide
	// GAZPtr map[PointerCstrPtrAlias[R, Z]]int // Soft reject

	// GACi     map[ComparableCstrAlias[int]]int                                         // Accept, hide
	// GACV     map[ComparableCstrAlias[ValueTextMarshaler]]int                          // Accept
	// GACP     map[ComparableCstrAlias[PointerTextMarshaler]]int                        // Reject
	// GACiPtr  map[ComparableCstrPtrAlias[int]]int                                      // Reject
	// GACVPtr  map[ComparableCstrPtrAlias[ValueTextMarshaler]]int                       // Accept, hide
	// GACPPtr  map[ComparableCstrPtrAlias[PointerTextMarshaler]]int                     // Accept, hide
	// GABi     map[BasicCstrAlias[int]]int                                              // Accept, hide
	// GABs     map[BasicCstrAlias[string]]int                                           // Accept
	// GABiPtr  map[BasicCstrPtrAlias[int]]int                                           // Reject
	// GABT     map[BadTildeCstrAlias[struct{}]]int                                      // Reject
	// GABTPtr  map[BadTildeCstrPtrAlias[struct{}]]int                                   // Reject
	// GAGT     map[GoodTildeCstrAlias[ValueTextMarshaler]]int                           // Accept
	// GAGTPtr  map[GoodTildeCstrPtrAlias[ValueTextMarshaler]]int                        // Accept, hide
	// GANBV    map[NonBasicCstrAlias[ValueTextMarshaler]]int                            // Accept
	// GANBP    map[NonBasicCstrAlias[*PointerTextMarshaler]]int                         // Accept, hide
	// GANBVPtr map[NonBasicCstrPtrAlias[ValueTextMarshaler]]int                         // Accept, hide
	// GANBPPtr map[NonBasicCstrPtrAlias[*PointerTextMarshaler]]int                      // Reject
	// GAPlV1   map[PointableCstrAlias[ValueTextMarshaler]]int                           // Accept
	// GAPlV2   map[*PointableCstrAlias[ValueTextMarshaler]]int                          // Accept
	// GAPlP1   map[PointableCstrAlias[PointerTextMarshaler]]int                         // Reject
	// GAPlP2   map[*PointableCstrAlias[PointerTextMarshaler]]int                        // Accept
	// GAPlVPtr map[PointableCstrPtrAlias[ValueTextMarshaler]]int                        // Accept, hide
	// GAPlPPtr map[PointableCstrPtrAlias[PointerTextMarshaler]]int                      // Accept, hide
	// GAMi     map[MixedCstrAlias[uint]]int                                             // Accept, hide
	// GAMS     map[MixedCstrAlias[StringType]]int                                       // Accept
	// GAMV     map[MixedCstrAlias[ValueTextMarshaler]]int                               // Accept
	// GAMSPtr  map[MixedCstrPtrAlias[StringType]]int                                    // Reject
	// GAMVPtr  map[MixedCstrPtrAlias[ValueTextMarshaler]]int                            // Accept, hide
	// GAII     map[InterfaceCstrAlias[encoding.TextMarshaler]]int                       // Accept, hide
	// GAIV     map[InterfaceCstrAlias[ValueTextMarshaler]]int                           // Accept
	// GAIP     map[InterfaceCstrAlias[*PointerTextMarshaler]]int                        // Accept, hide
	// GAIIPtr  map[InterfaceCstrPtrAlias[encoding.TextMarshaler]]int                    // Reject
	// GAIVPtr  map[InterfaceCstrPtrAlias[ValueTextMarshaler]]int                        // Accept, hide
	// GAIPPtr  map[InterfaceCstrPtrAlias[*PointerTextMarshaler]]int                     // Reject
	// GAPrV    map[PointerCstrAlias[ValueTextMarshaler, *ValueTextMarshaler]]int        // Accept, hide
	// GAPrP    map[PointerCstrAlias[PointerTextMarshaler, *PointerTextMarshaler]]int    // Accept, hide
	// GAPrVPtr map[PointerCstrPtrAlias[ValueTextMarshaler, *ValueTextMarshaler]]int     // Reject
	// GAPrPPtr map[PointerCstrPtrAlias[PointerTextMarshaler, *PointerTextMarshaler]]int // Reject
}

func (*Service) Method() (_ Maps[PointerTextMarshaler, int, int, ValueTextMarshaler, *PointerTextMarshaler, ValueTextMarshaler, StringType, ValueTextMarshaler, *PointerTextMarshaler]) {
	return
}

func main() {
	app := application.New(application.Options{
		Services: []application.Service{
			application.NewService(&Service{}),
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
