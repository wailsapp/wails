package idlreader

import (
	"os"
	"strings"
	"testing"
)

const fixture = `# synthetic declarations
import "base.idl";
[uuid(11111111-2222-3333-4444-555555555555),version(1.0)]
library Example {
 interface IExample;
 cpp_quote("example \"quoted\" text")
 [v1_enum] typedef enum Mode { First, Second=0x10, Third=1 << 3, } Mode;
 typedef struct Flags { UINT32 count; BOOL enabled; BYTE tag; } Flags;
 [uuid(aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee),object,pointer_default(unique)]
 interface IExample : IUnknown {
  [propget] HRESULT Value([out,retval] BOOL * value);
  [propput] HRESULT Value([in] BOOL value);
  HRESULT add_Changed([in] IHandler * handler, [out] Token * token);
  HRESULT remove_Changed([in] Token token);
  HRESULT Data([in,size_is(count)] BYTE const **data, [in] UINT32 count,);
 }
}`

func TestGrammar(t *testing.T) {
	for _, source := range []string{fixture, strings.ReplaceAll(fixture, "\n", "\r\n")} {
		idl, err := Parse([]byte(source))
		if err != nil {
			t.Fatal(err)
		}
		if len(idl.Imports) != 1 || idl.Imports[0].Name != `"base.idl"` {
			t.Fatalf("imports %+v", idl.Imports)
		}
		library := idl.Libraries[0]
		if library.Name != "Example" || len(library.Declarations) != 5 {
			t.Fatalf("library %+v", library)
		}
		enumeration := library.Declarations[2].Enum
		if enumeration.Values[0].Value != nil || enumeration.Values[1].Value.Value != "0x10" || *enumeration.Values[2].Value.LeftShift != 3 {
			t.Fatal("enum values changed")
		}
		iface := library.Declarations[4].Interface
		if *iface.Methods[0].Prop != "Get" || *iface.Methods[1].Prop != "Put" || iface.Methods[2].Name != "AddChanged" || iface.Methods[3].Name != "RemoveChanged" {
			t.Fatal("method naming changed")
		}
		param := iface.Methods[4].Params[0]
		if param.Direction.Retval != "size_is(count)" || param.Const != "const" || param.Pointer != "**" {
			t.Fatalf("parameter %+v", param)
		}
	}
}

func TestOctalShift(t *testing.T) {
	tree, err := Parse([]byte(strings.Replace(fixture, "Third=1 << 3", "Third=1 << 010", 1)))
	if err != nil {
		t.Fatal(err)
	}
	shift := tree.Libraries[0].Declarations[2].Enum.Values[2].Value.LeftShift
	if *shift != 8 {
		t.Fatalf("shift %d; want 8", *shift)
	}
}

func TestInvalidInput(t *testing.T) {
	for _, source := range []string{
		`import "unterminated`, "@", string([]byte{0}),
		strings.Replace(fixture, "version(1.0)", "version(nope)", 1),
		strings.Replace(fixture, "11111111-2222", "bad-uuid", 1),
		strings.Replace(fixture, "UINT32 count;", "UNKNOWN count;", 1),
		strings.Replace(fixture, "[propget]", "[unknown]", 1),
		strings.Replace(fixture, "size_is(count)", "size_is(1)", 1),
		strings.Replace(fixture, "[out,retval]", "[bad,retval]", 1),
		strings.Replace(fixture, "[out,retval]", "[out,other]", 1),
		strings.Replace(fixture, "[out,retval]", "[out,retval,other]", 1),
		strings.Replace(fixture, "Third=1 << 3", "Third=1 << 9999999999999999999999999", 1),
		strings.Replace(fixture, "Third=1 << 3", "Third=1 << 08", 1),
		fixture + " unexpected", strings.TrimSuffix(fixture, "}"),
	} {
		if _, err := Parse([]byte(source)); err == nil {
			t.Fatalf("accepted invalid input: %.80q", source)
		} else if !strings.Contains(err.Error(), "IDL:") {
			t.Fatalf("missing source location: %v", err)
		}
	}
	if tree, err := Parse(nil); err != nil || len(tree.Libraries) != 0 {
		t.Fatalf("empty IDL: %+v, %v", tree, err)
	}
}

func FuzzParse(f *testing.F) {
	f.Add([]byte(fixture))
	f.Add([]byte("# comment without newline"))
	f.Add([]byte("[uuid("))
	f.Fuzz(func(t *testing.T, data []byte) {
		tree, err := Parse(data)
		if err != nil {
			return
		}
		if tree == nil {
			t.Fatal("nil tree without error")
		}
		for _, library := range tree.Libraries {
			if library == nil || library.Header == nil || library.Name == "" {
				t.Fatal("incomplete library")
			}
		}
	})
}

func BenchmarkParse(b *testing.B) {
	data, err := os.ReadFile("../../WebView2.1.0.2903.40.idl")
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.SetBytes(int64(len(data)))
	for b.Loop() {
		if _, err := Parse(data); err != nil {
			b.Fatal(err)
		}
	}
}
