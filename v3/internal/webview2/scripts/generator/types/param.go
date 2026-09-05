package types

import (
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
	VtableCallInput string
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

func (p *Param) AsInputType() string {
	if p.isPointer() && p.GoType != "string" {
		return "*" + p.GoType
	}
	return p.GoType
}

func (p *Param) processSetup() {
	p.processSetupInputs()
	p.processSetupOutputs()
	p.processVtableCallInput()
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

func (p *Param) processVtableCallInput() {
	variableName := p.GetVariableName()
	if strings.HasPrefix(p.Type, "int") || strings.HasPrefix(p.Type, "uint") || p.Type == "bool" || p.Type == "float32" || p.Type == "float64" {
		p.VtableCallInput = "uintptr(" + variableName + ")"
		return
	}
	switch p.Type {
	case "LPCWSTR", "LPWSTR":
		p.VtableCallInput = "uintptr(unsafe.Pointer(" + variableName + "))"
		return
	}
	if p.Pointer == "**" {
		p.VtableCallInput = "uintptr(unsafe.Pointer(&" + variableName + "))"
		return
	}
	if p.Pointer == "*" {
		if p.IsOutputParam() {
			p.VtableCallInput = "uintptr(unsafe.Pointer(&" + variableName + "))"
		} else {
			p.VtableCallInput = "uintptr(unsafe.Pointer(" + variableName + "))"
		}
		return
	}
	if p.IsEnum() {
		p.VtableCallInput = "uintptr(" + variableName + ")"
		return
	}
	p.VtableCallInput = "uintptr(unsafe.Pointer(&" + variableName + "))"
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
	switch p.GoType {
	case "string":
		// We need to convert to *uint16
		p.setupTemplate = "inputStringSetup.tmpl"
		p.LocalName = "_" + p.Name
	}
}

func (p *Param) processSetupOutputs() {
	if !p.IsOutputParam() {
		return
	}
	switch p.GoType {
	case "string":
		p.LocalName = "_" + p.Name
		p.setupTemplate = "outputStringSetup.tmpl"
		p.cleanupTemplate = "outputStringCleanup.tmpl"
	case "bool":
		p.LocalName = "_" + p.Name
		p.setupTemplate = "outputBoolSetup.tmpl"
		p.cleanupTemplate = "outputBoolCleanup.tmpl"
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
