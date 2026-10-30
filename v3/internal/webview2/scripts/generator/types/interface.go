package types

import (
	"bytes"
	"io"
	"log"
	"slices"
	"strings"
	"text/template"
)

type InterfaceDeclaration struct {
	Header    *InterfaceHeader
	Name      string
	BaseClass string
	Methods   []*InterfaceMethod

	// private
	decl         *Declaration
	InvokeMethod *InterfaceMethod
	includes     []string
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
	d.addImport(`"unsafe"`)
	if len(d.Methods) == 1 && d.Methods[0] == d.InvokeMethod {
		return nil
	}
	d.addImport(`"syscall"`)
	d.addImport(`"golang.org/x/sys/windows"`)
	return nil
}

func (d *InterfaceDeclaration) addImport(path string) {
	if !slices.Contains(d.includes, path) {
		d.includes = append(d.includes, path)
	}
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
		Header          *InterfaceHeader
	}{
		PackageName:     packageName,
		BaseClass:       d.BaseClass,
		Header:          d.Header,
		Name:            d.Name,
		Methods:         d.Methods,
		HasInvokeMethod: d.HasInvokeMethod(),
		Includes:        d.includes,
	}
	if d.BaseClass == "IUnknown" {
		data.BaseClass = ""
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
	Prop       *Prop
	ReturnType string
	Name       InterfaceMethodName
	Params     []*Param

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
	m.processParams()

	return nil
}

func (m *InterfaceMethod) processParams() {
	for _, param := range m.Params {
		param.Process(m)
		if param.IsOutputParam() {
			m.outputParams = append(m.outputParams, param)
		} else {
			m.inputParams = append(m.inputParams, param)
		}
	}

	m.processInputParams()
	m.processOutputParams()
}

func (m *InterfaceMethod) processInputParams() {
	var inputs []string
	var inputParamNames []string
	for _, param := range m.inputParams {
		inputs = append(inputs, param.Name+" "+param.AsInputType())
		inputParamNames = append(inputParamNames, param.Name)
		param.processSetup()
	}
	m.GoInputs = strings.Join(inputs, ", ")
	m.InputParamNames = strings.Join(inputParamNames, ", ")
}

func (m *InterfaceMethod) processOutputParams() {
	var outputs []string
	var outputParamNames []string
	var outputParamTypes []string
	for _, param := range m.outputParams {
		outputs = append(outputs, param.Name+" "+param.GoType)
		outputParamNames = append(outputParamNames, param.Name)
		outputParamTypes = append(outputParamTypes, param.GoType)
		param.processSetup()
	}
	// Add the mandatory error
	outputs = append(outputs, "err error")
	outputParamNames = append(outputParamNames, "err")
	outputParamTypes = append(outputParamTypes, "error")

	m.GoOutputs = strings.Join(outputs, ", ")
	m.OutputParamNames = strings.Join(outputParamNames, ", ")
	m.GoReturnTypes = strings.Join(outputParamTypes, ", ")
	if len(outputParamTypes) > 1 {
		m.GoReturnTypes = "(" + m.GoReturnTypes + ")"
	}
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

func (m *InterfaceMethod) VtableCallInputs() string {
	var buffer bytes.Buffer
	for _, input := range m.Params {
		buffer.WriteString("\t\t" + input.VtableCallInput + ",\n")
	}
	return buffer.String()
}

func (m *InterfaceMethod) ReturnsHRESULT() bool {
	return m.ReturnType == "HRESULT"
}

func (m *InterfaceMethod) ErrorValues() string {
	var errorValues []string
	for _, outputParam := range m.outputParams {
		errorValues = append(errorValues, outputParam.defaultErrorValue())
	}
	errorValues = append(errorValues, "err")
	return strings.Join(errorValues, ", ")
}
func (m *InterfaceMethod) ErrorValuesHRESULT() string {
	var errorValues []string
	for _, outputParam := range m.outputParams {
		errorValues = append(errorValues, outputParam.defaultErrorValue())
	}
	if m.ReturnsHRESULT() {
		errorValues = append(errorValues, "syscall.Errno(hr)")
	} else {
		errorValues = append(errorValues, "err")
	}
	return strings.Join(errorValues, ", ")
}

func (m *InterfaceMethod) GetHResultVariable() string {
	if m.ReturnsHRESULT() {
		return "hr"
	}
	return "_"
}

func (m *InterfaceMethod) SuccessValues() string {
	var successValues []string
	for _, outputParam := range m.outputParams {
		successValues = append(successValues, outputParam.GetReturnVariableName())
	}
	successValues = append(successValues, "err")
	return strings.Join(successValues, ", ")
}

type InterfaceHeader struct {
	UUID *UUID
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
