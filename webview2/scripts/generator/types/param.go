package types

import (
	"fmt"
	"io"
	"strings"
)

type Param struct {
	Direction *Direction `parser:"@@?"`
	Type      string     `parser:"@Ident"`
	Const     string     `parser:"@('const')?"`
	Pointer   string     `parser:"@('*')*"`
	Name      string     `parser:"@Ident ','?"`

	// Processed
	GoType       string
	InputGoType  string
	OutputGoType string

	// This is used to generate setup code for the Go inputs
	setupTemplate   string
	cleanupTemplate string
	LocalName       string
	decl            *InterfaceMethod

	// callWords* hold the uintptr argument expression(s) this parameter
	// contributes to the vtable call, per architecture. Most parameters are a
	// single machine word everywhere; 8-byte values (double, [U]INT64, 8-byte
	// structs) are two stack slots on 386, and 16-byte structs are a register
	// pair on arm64, four stack words on 386 and a pointer-to-copy on amd64.
	callWordsAMD64 []string
	callWordsARM64 []string
	callWords386   []string
}

func (p *Param) IsOutputParam() bool {
	if p.Direction == nil {
		return false
	}
	return p.Direction.Dir == "out"
}

func (p *Param) LocalVariableType() string {
	//if p.IsOutputParam() && p.isDoublePointer() {
	//	return p.GoType[1:]
	//}
	return p.GoType
}

func (p *Param) Process(decl *InterfaceMethod) {
	p.decl = decl
	p.GoType = IdlTypeToGoType(p.Type)
	if p.isDoublePointer() {
		p.GoType = "*" + p.GoType
	}
	// LPCWSTR* (or LPWSTR*) in an [in] direction is a C array of strings, not a
	// single string with an extra level of indirection.
	if p.IsInputParam() && p.isSinglePointer() &&
		(p.Type == "LPCWSTR" || p.Type == "LPWSTR") {
		p.GoType = "[]string"
	}
	// [in] T** is a C array of interface pointers (paired with a count param);
	// surface it as []*T.
	if p.IsInputParam() && p.isDoublePointer() &&
		p.Type != "LPCWSTR" && p.Type != "LPWSTR" {
		p.GoType = "[]" + p.GoType
	}
	// [out] T*** — COM allocates an array of T* and writes its address back;
	// surface it as []*T (copied out, array freed with CoTaskMemFree).
	if p.IsOutputParam() && p.isTriplePointer() {
		p.GoType = "[]*" + IdlTypeToGoType(p.Type)
	}
	// [out] LPWSTR** — COM allocates an array of UTF-16 strings and writes
	// its address back; surface it as []string (each string decoded and
	// freed, then the array freed).
	if p.IsOutputParam() && p.isDoublePointer() &&
		(p.Type == "LPCWSTR" || p.Type == "LPWSTR") {
		p.GoType = "[]string"
	}
	p.OutputGoType = p.GoType
	if p.IsOutputParam() && strings.HasPrefix(p.OutputGoType, "**") {
		p.OutputGoType = p.GoType[1:]
	}
	p.InputGoType = p.GoType
}

func (p *Param) isPointer() bool {
	return p.Pointer != ""
}

func (p *Param) isSinglePointer() bool {
	return p.Pointer == "*"
}

func (p *Param) isDoublePointer() bool {
	return p.Pointer == "**"
}

func (p *Param) isTriplePointer() bool {
	return p.Pointer == "***"
}

func (p *Param) AsInputType() string {
	// Slice types ([]string, []*T) already encode the pointer semantics;
	// don't add another *.
	if strings.HasPrefix(p.GoType, "[]") {
		return p.GoType
	}
	if p.isPointer() && p.GoType != "string" {
		return "*" + p.GoType
	}
	return p.GoType
}

// ArrayElementType returns the element type of a slice-typed parameter,
// e.g. "*ICoreWebView2CustomSchemeRegistration" for a []*T array.
func (p *Param) ArrayElementType() string {
	return strings.TrimPrefix(p.GoType, "[]")
}

// CountVariableName returns the variable holding the element count that pairs
// with an array parameter: the first other output parameter of type uint32.
func (p *Param) CountVariableName() (string, error) {
	for _, op := range p.decl.outputParams {
		if op != p && op.GoType == "uint32" {
			return op.GetVariableName(), nil
		}
	}
	return "", fmt.Errorf("array parameter %q has no paired uint32 count output parameter", p.Name)
}

func (p *Param) processSetup() error {
	p.processSetupInputs()
	p.processSetupOutputs()
	return p.processVtableCallInput()
}

func (p *Param) SetupCode(w io.Writer) {
	if p.setupTemplate == "" {
		return
	}
	data := struct {
		Param       *Param
		ErrorValues string
	}{
		Param:       p,
		ErrorValues: p.decl.ErrorValues(),
	}
	mustTemplate("Param Setup: "+p.setupTemplate, p.setupTemplate, &data, w)
}
func (p *Param) CleanupCode(w io.Writer) {
	if p.cleanupTemplate == "" {
		return
	}
	mustTemplate("Param Cleanup: "+p.cleanupTemplate, p.cleanupTemplate, p, w)
}

func (p *Param) IsInputParam() bool {
	return !p.IsOutputParam()
}

// wordsAllArches sets the same single-word argument expression for every
// architecture.
func (p *Param) wordsAllArches(expr string) {
	p.callWordsAMD64 = []string{expr}
	p.callWordsARM64 = []string{expr}
	p.callWords386 = []string{expr}
}

// words64And386 sets a single-word expression for the 64-bit architectures and
// a low/high word pair for 386, where 8-byte values occupy two 4-byte stack
// slots (low word pushed at the lower address).
func (p *Param) words64And386(expr64, lo386, hi386 string) {
	p.callWordsAMD64 = []string{expr64}
	p.callWordsARM64 = []string{expr64}
	p.callWords386 = []string{lo386, hi386}
}

// builtinByValueSizes holds the sizes of imported Win32 types that appear by
// value in the IDL but are not declared as typedef structs inside it.
var builtinByValueSizes = map[string]int{
	"EventRegistrationToken": 8,
	"POINT":                  8,
	"RECT":                   16,
	// VARIANT deliberately absent: its size differs per architecture (24 bytes
	// on win64, 16 on win32) so it must never be passed by value. The IDL only
	// uses VARIANT*; if a by-value VARIANT ever appears, generation fails loudly.
}

func (p *Param) byValueSize() (int, bool) {
	if size, ok := builtinByValueSizes[p.Type]; ok {
		return size, true
	}
	size, ok := p.decl.decl.decl.library.structSizes[p.Type]
	return size, ok
}

func (p *Param) processVtableCallInput() error {
	variableName := p.GetVariableName()

	// String types: direction determines whether to pass pointer-to-pointer or pointer.
	// For output LPWSTR* the local var is *uint16; pass &local so COM writes the pointer back.
	// For input LPWSTR/LPCWSTR (plain or array) pass the *uint16 / **uint16 directly.
	switch p.Type {
	case "LPCWSTR", "LPWSTR":
		// Output [out] LPWSTR* needs &var so COM writes the *uint16 pointer back.
		// Input [in] LPWSTR* (array of strings, marshaled to **uint16 via
		// inputStringArraySetup.tmpl) and plain [in] LPWSTR (single *uint16)
		// both pass the local variable directly — the template chose the right Go type.
		if p.IsOutputParam() {
			p.wordsAllArches("uintptr(unsafe.Pointer(&" + variableName + "))")
		} else {
			p.wordsAllArches("uintptr(unsafe.Pointer(" + variableName + "))")
		}
		return nil
	}

	// Pointer checks come before the numeric GoType check so that output numeric
	// pointers (e.g. [out] int* / UINT32*) are correctly passed by address, not value.
	// Pointers are a single machine word on every architecture.
	if p.isTriplePointer() {
		if p.IsOutputParam() {
			// The local (from outputInterfaceArraySetup.tmpl) is a **T; COM
			// writes the array address through &local.
			p.wordsAllArches("uintptr(unsafe.Pointer(&" + variableName + "))")
			return nil
		}
		return fmt.Errorf("[in] triple pointer parameter %q of type %s is not supported", p.Name, p.Type)
	}
	if p.Pointer == "**" {
		if p.IsOutputParam() {
			// var local *T — COM writes the pointer through &local.
			p.wordsAllArches("uintptr(unsafe.Pointer(&" + variableName + "))")
		} else {
			// Input arrays: the local (from inputInterfaceArraySetup.tmpl) is
			// already a **T — pass it directly, NOT its address.
			p.wordsAllArches("uintptr(unsafe.Pointer(" + variableName + "))")
		}
		return nil
	}
	if p.Pointer == "*" {
		if p.IsOutputParam() {
			p.wordsAllArches("uintptr(unsafe.Pointer(&" + variableName + "))")
		} else {
			p.wordsAllArches("uintptr(unsafe.Pointer(" + variableName + "))")
		}
		return nil
	}

	if p.IsEnum() {
		// IDL enums are C ints: 4 bytes, one word everywhere.
		p.wordsAllArches("uintptr(" + variableName + ")")
		return nil
	}

	goType := p.GoType
	switch goType {
	case "float64":
		// Win64 syscalls mirror the first four integer registers into XMM0-3,
		// so passing the raw IEEE-754 bit pattern as an integer word reaches
		// the callee's double argument correctly on amd64; arm64 behaves the
		// same via its float mirroring. uintptr(v) would TRUNCATE the float to
		// its integer part. On 386 a double is two 4-byte stack slots.
		p.decl.decl.includes.AddUnique(`"math"`)
		p.words64And386(
			"uintptr(math.Float64bits("+variableName+"))",
			"uintptr(uint32(math.Float64bits("+variableName+")))",
			"uintptr(uint32(math.Float64bits("+variableName+")>>32))",
		)
		return nil
	case "float32":
		p.decl.decl.includes.AddUnique(`"math"`)
		p.wordsAllArches("uintptr(math.Float32bits(" + variableName + "))")
		return nil
	case "int64", "uint64":
		p.words64And386(
			"uintptr("+variableName+")",
			"uintptr(uint32(uint64("+variableName+")))",
			"uintptr(uint32(uint64("+variableName+")>>32))",
		)
		return nil
	}

	// Scalar numeric / bool inputs: use GoType (handles uppercase IDL names like UINT32,
	// INT32, BOOL that map to Go uint32, int32, bool). For bool, setup code converts to
	// int32 first so the local variable is already an int32 and uintptr() is correct.
	//
	// HWND, HANDLE, HCURSOR are uintptr aliases in golang.org/x/sys/windows — they're
	// already handle *values*, not pointers to handles. Passing their address gives a
	// stack-local pointer instead of the handle itself, which the receiving COM method
	// dereferences as garbage. Match the existing defaultErrorValue() treatment of
	// these as scalars.
	if strings.HasPrefix(goType, "int") || strings.HasPrefix(goType, "uint") ||
		goType == "bool" ||
		goType == "HWND" || goType == "HANDLE" || goType == "HCURSOR" {
		p.wordsAllArches("uintptr(" + variableName + ")")
		return nil
	}

	// Everything left is an aggregate passed by value. The Win64 ABI passes
	// 1/2/4/8-byte aggregates by value in a register — NOT by address. The
	// ARM64 ABI additionally passes 9-16 byte composites in a register pair.
	// 386 stdcall pushes the whole aggregate onto the stack as 4-byte words.
	size, ok := p.byValueSize()
	if !ok {
		return fmt.Errorf("don't know how to marshal by-value parameter %q of type %s — add the type to builtinByValueSizes or declare its struct in the IDL", p.Name, p.Type)
	}
	switch size {
	case 1, 2, 4:
		load := map[int]string{1: "uint8", 2: "uint16", 4: "uint32"}[size]
		p.wordsAllArches("uintptr(*(*" + load + ")(unsafe.Pointer(&" + variableName + ")))")
	case 8:
		p.words64And386(
			"uintptr(*(*uint64)(unsafe.Pointer(&"+variableName+")))",
			"uintptr((*(*[2]uint32)(unsafe.Pointer(&"+variableName+")))[0])",
			"uintptr((*(*[2]uint32)(unsafe.Pointer(&"+variableName+")))[1])",
		)
	case 16:
		// amd64 passes >8-byte aggregates via pointer to a copy; the Go
		// parameter is already a private copy so its address is safe to hand
		// over. arm64 packs 9-16 byte composites into two registers. 386
		// pushes four stack words.
		p.callWordsAMD64 = []string{"uintptr(unsafe.Pointer(&" + variableName + "))"}
		p.callWordsARM64 = []string{
			"(*(*[2]uintptr)(unsafe.Pointer(&" + variableName + ")))[0]",
			"(*(*[2]uintptr)(unsafe.Pointer(&" + variableName + ")))[1]",
		}
		words := make([]string, 4)
		for i := range words {
			words[i] = fmt.Sprintf("uintptr((*(*[4]uint32)(unsafe.Pointer(&%s)))[%d])", variableName, i)
		}
		p.callWords386 = words
	default:
		return fmt.Errorf("by-value parameter %q of type %s has unsupported size %d — extend processVtableCallInput", p.Name, p.Type, size)
	}
	return nil
}

func (p *Param) ClearLocalName() string {
	p.LocalName = ""
	return ""
}

func (p *Param) GetVariableName() string {
	result := p.LocalName
	if result == "" {
		result = p.Name
	}
	return result
}

func (p *Param) GetReturnVariableName() string {
	result := p.LocalName
	if result == "" {
		result = p.Name
	}
	return result
}

func (p *Param) IsEnum() bool {
	return p.decl.decl.decl.library.enums.Contains(p.Type)
}

func (p *Param) processSetupInputs() {
	if !p.IsInputParam() {
		return
	}
	switch {
	case p.GoType == "string":
		p.setupTemplate = "inputStringSetup.tmpl"
		p.LocalName = "_" + p.Name
	case p.GoType == "[]string":
		// LPCWSTR* — convert Go slice to a C array of *uint16 pointers
		p.setupTemplate = "inputStringArraySetup.tmpl"
		p.LocalName = "_" + p.Name
	case strings.HasPrefix(p.GoType, "[]*"):
		// T** — convert Go slice to a pointer to its first element
		p.setupTemplate = "inputInterfaceArraySetup.tmpl"
		p.LocalName = "_" + p.Name
	case p.GoType == "bool":
		// COM BOOL is int32; convert before the vtable call
		p.setupTemplate = "inputBoolSetup.tmpl"
		p.LocalName = "_" + p.Name
	}
}

func (p *Param) processSetupOutputs() {
	if !p.IsOutputParam() {
		return
	}
	switch {
	case p.GoType == "string":
		p.LocalName = "_" + p.Name
		p.setupTemplate = "outputStringSetup.tmpl"
		p.cleanupTemplate = "outputStringCleanup.tmpl"
	case p.GoType == "bool":
		p.LocalName = "_" + p.Name
		p.setupTemplate = "outputBoolSetup.tmpl"
		p.cleanupTemplate = "outputBoolCleanup.tmpl"
	case strings.HasPrefix(p.GoType, "[]*"):
		// [out] T*** — COM-allocated array copied into a Go slice
		p.LocalName = "_" + p.Name
		p.setupTemplate = "outputInterfaceArraySetup.tmpl"
		p.cleanupTemplate = "outputInterfaceArrayCleanup.tmpl"
	case p.GoType == "[]string":
		// [out] LPWSTR** — COM-allocated string array decoded and freed
		p.LocalName = "_" + p.Name
		p.setupTemplate = "outputStringArraySetup.tmpl"
		p.cleanupTemplate = "outputStringArrayCleanup.tmpl"
	default:
		p.setupTemplate = "outputDefaultSetup.tmpl"
	}
	if p.Pointer != "" {
		p.decl.decl.includes.AddUnique(`"unsafe"`)
	}
}

func (p *Param) defaultErrorValue() string {
	switch true {
	case p.IsEnum(), strings.HasPrefix(p.GoType, "uint"), strings.HasPrefix(p.GoType, "int"),
		p.GoType == "HANDLE", p.GoType == "HWND", p.GoType == "HCURSOR":
		return "0"
	case strings.HasPrefix(p.GoType, "float"):
		return "0.0"
	case p.GoType == "bool":
		return "false"
	case p.GoType == "string":
		return `""`
	case strings.HasPrefix(p.GoType, "[]"):
		return "nil"
	case p.OutputGoType[0] == '*':
		return "nil"
	default:
		return p.GoType + "{}"
	}
}

type Direction struct {
	Dir    string `parser:"'[' @('out'|'in')"`
	Retval string `parser:"(',' @('retval'|'size_is' '(' Ident ')') )? ']'"`
}
