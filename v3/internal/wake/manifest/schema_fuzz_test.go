package manifest

import (
	"testing"
	"unicode/utf8"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func FuzzManifestDecodeNeverPanics(f *testing.F) {
	f.Add([]byte(""))
	f.Add([]byte("version = 3"))
	f.Add([]byte("version = 3\nproject {}"))
	f.Add([]byte("version = 3\nproject { name = var.name }"))
	f.Add([]byte("version = 3\nproject { name = \"\u2603\" }"))
	root := f.TempDir()
	f.Fuzz(func(t *testing.T, source []byte) {
		_, _ = decodeHCL(root, "fuzz.hcl", source, "")
	})
}

func FuzzManifestLiteralStringRoundTrip(f *testing.F) {
	f.Add("")
	f.Add("plain")
	f.Add("quote: \" and slash: \\")
	f.Add("line one\nline two")
	f.Add("snowman: ☃")
	f.Add("control: \x00\x01")
	f.Add("carriage\rreturn\tand tab")
	f.Add("template markers: ${value} and %{if true}")
	root := f.TempDir()
	f.Fuzz(func(t *testing.T, literal string) {
		if !utf8.ValidString(literal) {
			t.Skip()
		}
		normalized := cty.StringVal(literal).AsString()
		literalToken := hclwrite.TokensForValue(cty.StringVal(normalized)).Bytes()
		source := []byte("version = 3\nproject {\n" +
			"  name = \"literal\"\n" +
			"  product_name = \"Literal\"\n" +
			"  identifier = \"com.example.literal\"\n" +
			"  version = \"1.0.0\"\n" +
			"  description = " + string(literalToken) + "\n}\n")
		loaded, err := decodeHCL(root, "fuzz.hcl", source, "")
		if err != nil {
			t.Fatalf("decode generated literal %q: %v", literal, err)
		}
		encoded, err := EncodeConfig(loaded.Config)
		if err != nil {
			t.Fatal(err)
		}
		roundTripped, err := decodeHCL(root, "round-trip.hcl", encoded, "")
		if err != nil {
			t.Fatal(err)
		}
		if got := roundTripped.Config.Project.Description; got != normalized {
			t.Fatalf("description = %q, want NFC-normalized %q", got, normalized)
		}
	})
}
