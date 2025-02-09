package main

import (
	_ "embed"
	"encoding"
	"encoding/json"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct{}

// class {}
type NonMarshaler struct{}

// any
type ValueJsonMarshaler struct{}

func (ValueJsonMarshaler) MarshalJSON() ([]byte, error) { return nil, nil }

// any
type PointerJsonMarshaler struct{}

func (*PointerJsonMarshaler) MarshalJSON() ([]byte, error) { return nil, nil }

// string
type ValueTextMarshaler struct{}

func (ValueTextMarshaler) MarshalText() ([]byte, error) { return nil, nil }

// string
type PointerTextMarshaler struct{}

func (*PointerTextMarshaler) MarshalText() ([]byte, error) { return nil, nil }

// any
type ValueMarshaler struct{}

func (ValueMarshaler) MarshalJSON() ([]byte, error) { return nil, nil }
func (ValueMarshaler) MarshalText() ([]byte, error) { return nil, nil }

// any
type PointerMarshaler struct{}

func (*PointerMarshaler) MarshalJSON() ([]byte, error) { return nil, nil }
func (*PointerMarshaler) MarshalText() ([]byte, error) { return nil, nil }

// any
type UnderlyingJsonMarshaler struct{ json.Marshaler }

// string
type UnderlyingTextMarshaler struct{ encoding.TextMarshaler }

// any
type UnderlyingMarshaler struct {
	json.Marshaler
	encoding.TextMarshaler
}

type customJsonMarshaler interface {
	MarshalJSON() ([]byte, error)
}

type customTextMarshaler interface {
	MarshalText() ([]byte, error)
}

type customMarshaler interface {
	MarshalJSON() ([]byte, error)
	MarshalText() ([]byte, error)
}

// struct{}
type AliasNonMarshaler = struct{}

// any
type AliasJsonMarshaler = struct{ json.Marshaler }

// string
type AliasTextMarshaler = struct{ encoding.TextMarshaler }

// any
type AliasMarshaler = struct {
	json.Marshaler
	encoding.TextMarshaler
}

// any
type ImplicitJsonMarshaler UnderlyingJsonMarshaler

// string
type ImplicitTextMarshaler UnderlyingTextMarshaler

// any
type ImplicitMarshaler UnderlyingMarshaler

// string
type ImplicitNonJson UnderlyingMarshaler

func (ImplicitNonJson) MarshalJSON() {}

// any
type ImplicitNonText UnderlyingMarshaler

func (ImplicitNonText) MarshalText() {}

// class{ Marshaler, TextMarshaler }
type ImplicitNonMarshaler UnderlyingMarshaler

func (ImplicitNonMarshaler) MarshalJSON() {}
func (ImplicitNonMarshaler) MarshalText() {}

// any
type ImplicitJsonButText UnderlyingJsonMarshaler

func (ImplicitJsonButText) MarshalText() ([]byte, error) { return nil, nil }

// any
type ImplicitTextButJson UnderlyingTextMarshaler

func (ImplicitTextButJson) MarshalJSON() ([]byte, error) { return nil, nil }

type Data struct {
	NM    NonMarshaler
	NMPtr *NonMarshaler // NonMarshaler | null

	VJM    ValueJsonMarshaler
	VJMPtr *ValueJsonMarshaler // ValueJsonMarshaler | null
	PJM    PointerJsonMarshaler
	PJMPtr *PointerJsonMarshaler // PointerJsonMarshaler | null

	VTM    ValueTextMarshaler
	VTMPtr *ValueTextMarshaler // ValueTextMarshaler | null
	PTM    PointerTextMarshaler
	PTMPtr *PointerTextMarshaler // PointerTextMarshaler | null

	VM    ValueMarshaler
	VMPtr *ValueMarshaler // ValueMarshaler | null
	PM    PointerMarshaler
	PMPtr *PointerMarshaler // PointerMarshaler | null

	UJM    UnderlyingJsonMarshaler
	UJMPtr *UnderlyingJsonMarshaler // UnderlyingJsonMarshaler | null
	UTM    UnderlyingTextMarshaler
	UTMPtr *UnderlyingTextMarshaler // UnderlyingTextMarshaler | null
	UM     UnderlyingMarshaler
	UMPtr  *UnderlyingMarshaler // UnderlyingMarshaler | null

	JM     struct{ json.Marshaler }          // any
	JMPtr  *struct{ json.Marshaler }         // any | null
	TM     struct{ encoding.TextMarshaler }  // string
	TMPtr  *struct{ encoding.TextMarshaler } // string | null
	CJM    struct{ customJsonMarshaler }     // any
	CJMPtr *struct{ customJsonMarshaler }    // any | null
	CTM    struct{ customTextMarshaler }     // string
	CTMPtr *struct{ customTextMarshaler }    // string | null
	CM     struct{ customMarshaler }         // any
	CMPtr  *struct{ customMarshaler }        // any | null

	ANM    AliasNonMarshaler
	ANMPtr *AliasNonMarshaler // AliasNonMarshaler | null
	AJM    AliasJsonMarshaler
	AJMPtr *AliasJsonMarshaler // AliasJsonMarshaler | null
	ATM    AliasTextMarshaler
	ATMPtr *AliasTextMarshaler // AliasTextMarshaler | null
	AM     AliasMarshaler
	AMPtr  *AliasMarshaler // AliasMarshaler | null

	ImJM    ImplicitJsonMarshaler
	ImJMPtr *ImplicitJsonMarshaler // ImplicitJsonMarshaler | null
	ImTM    ImplicitTextMarshaler
	ImTMPtr *ImplicitTextMarshaler // ImplicitTextMarshaler | null
	ImM     ImplicitMarshaler
	ImMPtr  *ImplicitMarshaler // ImplicitMarshaler | null

	ImNJ    ImplicitNonJson
	ImNJPtr *ImplicitNonJson // ImplicitNonJson | null
	ImNT    ImplicitNonText
	ImNTPtr *ImplicitNonText // ImplicitNonText | null
	ImNM    ImplicitNonMarshaler
	ImNMPtr *ImplicitNonMarshaler // ImplicitNonMarshaler | null

	ImJbT    ImplicitJsonButText
	ImJbTPtr *ImplicitJsonButText // ImplicitJsonButText | null
	ImTbJ    ImplicitTextButJson
	ImTbJPtr *ImplicitTextButJson // ImplicitTextButJson | null
}

func (*Service) Method() (_ Data) {
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
