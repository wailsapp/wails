package types

import (
	"bytes"
	"fmt"
	"github.com/leaanthony/slicer"
	"io"
	"log"
	"slices"
	"strings"
	"text/template"
)

type InterfaceDeclaration struct {
	Header    *InterfaceHeader   `parser:"'[' @@ ']'"`
	Name      string             `parser:"'interface' @Ident"`
	BaseClass string             `parser:" ':' @Ident '{' "`
	Methods   []*InterfaceMethod `parser:"@@+ '}'"`

	// private
	decl         *Declaration
	InvokeMethod *InterfaceMethod
	includes     slicer.StringSlicer
}

func (d *InterfaceDeclaration) Process(decl *Declaration) error {
	d.decl = decl

	// Find Invoke method
	for _, method := range d.Methods {
		err := method.Process(d)
		if err != nil {
			return err
		}
		if string(method.Name) == "Invoke" {
			d.InvokeMethod = method
			break
		}
	}
	d.includes.AddUnique(`"unsafe"`)
	if len(d.Methods) == 1 && d.Methods[0] == d.InvokeMethod {
		return nil
	}
	d.includes.AddUnique(`"syscall"`)
	d.includes.AddUnique(`"golang.org/x/sys/windows"`)
	return nil
}

func (d *InterfaceDeclaration) Generate(packageName string, w io.Writer) error {
	err := d.generateVtbl(packageName, w)
	if err != nil {
		return err
	}

	err = d.generateInvoke(w)
	if err != nil {
		return err
	}

	err = d.generateInterfaceMethods(w)
	if err != nil {
		return err
	}

	return nil
}

func (d *InterfaceDeclaration) generateVtbl(packageName string, w io.Writer) error {
	data := struct {
		PackageName     string
		Name            string
		Methods         []*InterfaceMethod
		HasInvokeMethod bool
		Includes        []string
		BaseClass       string
		BaseVtbl        string
		QIReceiver      string
		Header          *InterfaceHeader
	}{
		PackageName:     packageName,
		BaseClass:       d.BaseClass,
		BaseVtbl:        "IUnknownVtbl",
		Header:          d.Header,
		Name:            d.Name,
		Methods:         d.Methods,
		HasInvokeMethod: d.HasInvokeMethod(),
		Includes:        d.includes.AsSlice(),
	}
	library := d.decl.library
	if d.BaseClass == "IUnknown" {
		data.BaseClass = ""
	} else {
		// COM vtbls are flat: a derived interface's table starts with every
		// slot of its base. Embedding the base's vtbl struct (recursively
		// down to IUnknownVtbl) reproduces that layout; embedding only
		// IUnknownVtbl would shift every method onto the wrong slot.
		if _, ok := library.interfaceBases[d.BaseClass]; !ok {
			return fmt.Errorf("interface %s derives from %s which is not declared in the IDL", d.Name, d.BaseClass)
		}
		data.BaseVtbl = d.BaseClass + "Vtbl"
		// QueryInterface helpers must run against the object that actually
		// implements the interface: the root of the inheritance chain (e.g.
		// ICoreWebView2Controller for ICoreWebView2Controller2), not
		// ICoreWebView2 unconditionally.
		root, err := library.chainRoot(d.Name)
		if err != nil {
			return err
		}
		data.QIReceiver = root
	}
	mustTemplate("Interface Vtbl", "interfacevtbl.tmpl", &data, w)
	return nil
}

func (d *InterfaceDeclaration) GetBaseClass() string {
	if d.BaseClass == "IUnknown" {
		return ""
	}
	return d.BaseClass
}

func (d *InterfaceDeclaration) generateInvoke(w io.Writer) error {
	if d.InvokeMethod == nil {
		return nil
	}
	data := struct {
		Name         string
		InvokeMethod *InterfaceMethod
		Declaration  *InterfaceDeclaration
	}{
		Declaration:  d,
		Name:         d.Name,
		InvokeMethod: d.InvokeMethod,
	}
	mustTemplate("Interface Invoke", "interfaceInvoke.tmpl", &data, w)
	return nil
}

func (d *InterfaceDeclaration) HasInvokeMethod() bool {
	return d.InvokeMethod != nil
}

func mustTemplate(templateName string, filename string, data interface{}, w io.Writer) {
	templateData, err := templates.ReadFile("templates/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err := template.New(templateName).Parse(string(templateData))
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *InterfaceDeclaration) generateInterfaceMethods(w io.Writer) error {
	if len(d.Methods) == 1 && d.Methods[0] == d.InvokeMethod {
		return nil
	}
	for _, method := range d.Methods {
		data := struct {
			Name   string
			Method *InterfaceMethod
		}{
			Name:   d.Name,
			Method: method,
		}
		mustTemplate("Interface Methods", "interfaceMethod.tmpl", &data, w)
	}
	return nil
}

type InterfaceMethod struct {
	Prop       *Prop               `parser:"('[' @('propget'|'propput') ']')?"`
	ReturnType string              `parser:"@Ident"`
	Name       InterfaceMethodName `parser:"@Ident '('"`
	Params     []*Param            `parser:" @@* ')' ';'"`

	// private
	GoMethodName string

	GoInputs        string
	InputParamNames string

	GoReturnTypes string

	ProcessedName    string
	inputParams      []*Param
	outputParams     []*Param
	OutputParamNames string
	GoOutputs        string
	decl             *InterfaceDeclaration
}

func (m *InterfaceMethod) Process(decl *InterfaceDeclaration) error {
	m.decl = decl
	// Generate Go Method name
	goMethodName := strings.TrimPrefix(decl.Name, "ICoreWebView2")
	goMethodName = strings.TrimSuffix(goMethodName, "Handler")
	goMethodName = strings.TrimSuffix(goMethodName, "Event")
	m.GoMethodName = goMethodName

	m.ProcessedName = string(m.Name)
	if m.Prop != nil {
		m.ProcessedName = string(*m.Prop) + m.ProcessedName
	}
	return m.processParams()
}

func (m *InterfaceMethod) processParams() error {
	for _, param := range m.Params {
		param.Process(m)
		if param.IsOutputParam() {
			m.outputParams = append(m.outputParams, param)
		} else {
			m.inputParams = append(m.inputParams, param)
		}
	}

	if err := m.processInputParams(); err != nil {
		return fmt.Errorf("%s.%s: %w", m.decl.Name, m.ProcessedName, err)
	}
	if err := m.processOutputParams(); err != nil {
		return fmt.Errorf("%s.%s: %w", m.decl.Name, m.ProcessedName, err)
	}
	return nil
}

func (m *InterfaceMethod) processInputParams() error {
	var inputs slicer.StringSlicer
	var inputParamNames slicer.StringSlicer
	for _, param := range m.inputParams {
		inputs.Add(param.Name + " " + param.AsInputType())
		inputParamNames.Add(param.Name)
		if err := param.processSetup(); err != nil {
			return err
		}
	}
	m.GoInputs = inputs.Join(", ")
	m.InputParamNames = inputParamNames.Join(", ")
	return nil
}

func (m *InterfaceMethod) processOutputParams() error {
	var outputs slicer.StringSlicer
	var outputParamNames slicer.StringSlicer
	var outputParamTypes slicer.StringSlicer
	for _, param := range m.outputParams {
		outputs.Add(param.Name + " " + param.GoType)
		outputParamNames.Add(param.Name)
		outputParamTypes.Add(param.GoType)
		if err := param.processSetup(); err != nil {
			return err
		}
	}
	// Add the mandatory error
	outputs.Add("err error")
	outputParamNames.Add("err")
	outputParamTypes.Add("error")

	m.GoOutputs = outputs.Join(", ")
	m.OutputParamNames = outputParamNames.Join(", ")
	m.GoReturnTypes = outputParamTypes.Join(", ")
	if outputParamTypes.Length() > 1 {
		m.GoReturnTypes = "(" + m.GoReturnTypes + ")"
	}
	return nil
}

func (m *InterfaceMethod) SetupCode() string {
	var buffer bytes.Buffer
	for _, param := range m.Params {
		param.SetupCode(&buffer)
	}
	return buffer.String()
}

func (m *InterfaceMethod) CleanupCode() string {
	var buffer bytes.Buffer
	for _, param := range m.Params {
		param.CleanupCode(&buffer)
	}
	return buffer.String()
}

// NeedsArchSplit reports whether any parameter marshals differently across
// architectures (8-byte values on 386, 16-byte aggregates everywhere), which
// requires emitting one vtable call per architecture group.
func (m *InterfaceMethod) NeedsArchSplit() bool {
	for _, p := range m.Params {
		if !slices.Equal(p.callWordsAMD64, p.callWordsARM64) ||
			!slices.Equal(p.callWordsAMD64, p.callWords386) {
			return true
		}
	}
	return false
}

// NeedsARM64Variant reports whether arm64 marshalling differs from amd64
// (only true for 9-16 byte aggregates, which arm64 passes in a register pair
// where amd64 passes a pointer to a copy). When false, the arm64 case is
// folded into the default branch of the generated switch.
func (m *InterfaceMethod) NeedsARM64Variant() bool {
	for _, p := range m.Params {
		if !slices.Equal(p.callWordsAMD64, p.callWordsARM64) {
			return true
		}
	}
	return false
}

func (m *InterfaceMethod) joinCallWords(indent string, words func(*Param) []string) string {
	var buffer bytes.Buffer
	for _, input := range m.Params {
		for _, word := range words(input) {
			buffer.WriteString(indent + word + ",\n")
		}
	}
	return buffer.String()
}

func (m *InterfaceMethod) VtableCallInputs() string {
	return m.joinCallWords("\t\t", func(p *Param) []string { return p.callWordsAMD64 })
}

func (m *InterfaceMethod) VtableCallInputsAMD64() string {
	return m.joinCallWords("\t\t\t", func(p *Param) []string { return p.callWordsAMD64 })
}

func (m *InterfaceMethod) VtableCallInputsARM64() string {
	return m.joinCallWords("\t\t\t", func(p *Param) []string { return p.callWordsARM64 })
}

func (m *InterfaceMethod) VtableCallInputs386() string {
	return m.joinCallWords("\t\t\t", func(p *Param) []string { return p.callWords386 })
}

func (m *InterfaceMethod) ReturnsHRESULT() bool {
	return m.ReturnType == "HRESULT"
}

func (m *InterfaceMethod) ErrorValues() string {
	var errorValues slicer.StringSlicer
	for _, outputParam := range m.outputParams {
		errorValues.Add(outputParam.defaultErrorValue())
	}
	errorValues.Add("err")
	return errorValues.Join(", ")
}
func (m *InterfaceMethod) ErrorValuesHRESULT() string {
	var errorValues slicer.StringSlicer
	for _, outputParam := range m.outputParams {
		errorValues.Add(outputParam.defaultErrorValue())
	}
	if m.ReturnsHRESULT() {
		errorValues.Add("syscall.Errno(hr)")
	} else {
		errorValues.Add("err")
	}
	return errorValues.Join(", ")
}

func (m *InterfaceMethod) GetHResultVariable() string {
	if m.ReturnsHRESULT() {
		return "hr"
	}
	return "_"
}

// InvokeGoInputs renders the parameter list for the C-side trampoline that
// windows.NewCallback wraps. Strings are passed as *uint16 because
// windows.NewCallback panics at init if any parameter is wider than a uintptr
// (Go strings are 2-word fat pointers and slices are 3-word).
func (m *InterfaceMethod) InvokeGoInputs() string {
	var inputs slicer.StringSlicer
	for _, p := range m.inputParams {
		t := p.AsInputType()
		if t == "string" {
			t = "*uint16"
		}
		inputs.Add(p.Name + " " + t)
	}
	return inputs.Join(", ")
}

// InvokeConversionCode emits the Go statements that convert *uint16 trampoline
// parameters back into Go strings before the impl call. Empty if the Invoke
// method has no string parameters.
func (m *InterfaceMethod) InvokeConversionCode() string {
	var buf bytes.Buffer
	for _, p := range m.inputParams {
		if p.AsInputType() == "string" {
			buf.WriteString("\t_")
			buf.WriteString(p.Name)
			buf.WriteString(" := UTF16PtrToString(")
			buf.WriteString(p.Name)
			buf.WriteString(")\n")
		}
	}
	return buf.String()
}

// InvokeInputParamNames is the InputParamNames variant for the Invoke
// trampoline — string params become `_name` (the converted Go string),
// everything else stays as `name`.
func (m *InterfaceMethod) InvokeInputParamNames() string {
	var names slicer.StringSlicer
	for _, p := range m.inputParams {
		name := p.Name
		if p.AsInputType() == "string" {
			name = "_" + name
		}
		names.Add(name)
	}
	return names.Join(", ")
}

func (m *InterfaceMethod) SuccessValues() string {
	// The third return from syscall.Call is GetLastError, which is non-nil after
	// every call regardless of HRESULT — using it as the method's err return causes
	// successful calls to surface stale Win32 errors from prior unrelated syscalls.
	// The template binds the third return to `_`; the success path returns nil.
	var successValues slicer.StringSlicer
	for _, outputParam := range m.outputParams {
		successValues.Add(outputParam.GetReturnVariableName())
	}
	successValues.Add("nil")
	return successValues.Join(", ")
}

type InterfaceHeader struct {
	UUID *UUID `parser:"'uuid' '(' @UUID ')' ',' 'object' ',' 'pointer_default' '(' 'unique' ')'"`
}

func (h *InterfaceHeader) AsString() string {
	uuid := *h.UUID
	return string(`"{` + uuid + `}"`)
}

type InterfaceMethodName string

func (m *InterfaceMethodName) Capture(values []string) error {
	if len(values) == 0 {
		return nil
	}
	result := values[0]
	if strings.HasPrefix(values[0], "add_") {
		result = "Add" + result[4:]
	}
	if strings.HasPrefix(values[0], "remove_") {
		result = "Remove" + result[7:]
	}
	*m = InterfaceMethodName(result)
	return nil
}

type Prop string

func (p *Prop) Capture(values []string) error {
	if len(values) == 0 {
		return nil
	}
	result := strings.Title(values[0][4:])
	*p = Prop(result)
	return nil
}
