// Package idlreader reads the WebView2 IDL grammar consumed by the generator.
package idlreader

import (
	"fmt"
	"strconv"
	"strings"
	"updater/generator/types"
)

type parser struct {
	lexer   lexer
	current token
	err     error
}

func Parse(data []byte) (*types.IDL, error) {
	p := &parser{lexer: lexer{input: string(data), line: 1, column: 1}}
	p.next()
	result := &types.IDL{}
	for p.err == nil && p.current.text == "import" {
		p.next()
		var name strings.Builder
		for p.err == nil && p.current.text != ";" && p.current.kind != 0 {
			name.WriteString(p.current.text)
			p.next()
		}
		p.expect(";")
		result.Imports = append(result.Imports, &types.Import{Name: name.String()})
	}
	for p.err == nil && p.current.kind != 0 {
		result.Libraries = append(result.Libraries, p.library())
	}
	if p.err != nil {
		return nil, p.err
	}
	return result, nil
}

func (p *parser) next() {
	if p.err != nil {
		return
	}
	p.current, p.err = p.lexer.next()
}
func (p *parser) fail(message string) {
	if p.err == nil {
		p.err = fmt.Errorf("IDL:%d:%d: %s (got %q)", p.current.line, p.current.column, message, p.current.text)
	}
}
func (p *parser) take(text string) bool {
	if p.err != nil || p.current.text != text {
		return false
	}
	p.next()
	return true
}
func (p *parser) expect(text string) {
	if !p.take(text) {
		p.fail("expected " + strconv.Quote(text))
	}
}
func (p *parser) value(kind byte, label string) string {
	if p.err != nil {
		return ""
	}
	if p.current.kind != kind {
		p.fail("expected " + label)
		return ""
	}
	v := p.current.text
	p.next()
	return v
}
func (p *parser) ident() string { return p.value('i', "identifier") }

type attribute struct {
	name string
	args []token
}

func (p *parser) attributes() []attribute {
	p.expect("[")
	var result []attribute
	for p.err == nil {
		a := attribute{name: p.ident()}
		if p.take("(") {
			for p.err == nil && p.current.text != ")" && p.current.kind != 0 {
				a.args = append(a.args, p.current)
				p.next()
			}
			p.expect(")")
		}
		result = append(result, a)
		if !p.take(",") {
			break
		}
	}
	p.expect("]")
	return result
}
func (p *parser) attribute(a attribute, name string, args ...string) {
	if a.name != name || len(a.args) != len(args) {
		p.fail("expected attribute " + name)
		return
	}
	for i, value := range args {
		if a.args[i].text != value {
			p.fail("invalid arguments to " + name)
			return
		}
	}
}
func (p *parser) uuid(a attribute) string {
	if a.name != "uuid" || len(a.args) != 1 || a.args[0].kind != 'u' {
		p.fail("expected uuid attribute")
		return ""
	}
	return a.args[0].text
}
func (p *parser) library() *types.Library {
	attributes := p.attributes()
	l := &types.Library{Header: &types.LibraryHeader{}}
	if len(attributes) != 2 {
		p.fail("expected library uuid and version")
		return l
	}
	l.Header.UUID = p.uuid(attributes[0])
	version := attributes[1]
	if version.name != "version" || !(len(version.args) == 1 && version.args[0].kind == 'n' || len(version.args) == 3 && version.args[0].kind == 'n' && version.args[1].text == "." && version.args[2].kind == 'n') {
		p.fail("expected version attribute")
	}
	p.expect("library")
	l.Name = p.ident()
	p.expect("{")
	for p.err == nil && p.current.text != "}" {
		l.Declarations = append(l.Declarations, p.declaration())
	}
	p.expect("}")
	return l
}
func (p *parser) declaration() *types.Declaration {
	d := &types.Declaration{}
	switch p.current.text {
	case "interface":
		p.next()
		d.InterfaceForewardDecl = p.ident()
		p.expect(";")
	case "typedef":
		d.Struct = p.structure()
	case "cpp_quote":
		p.next()
		p.expect("(")
		d.CppQuote = p.value('s', "quoted C++ text")
		p.expect(")")
	case "[":
		attributes := p.attributes()
		if len(attributes) == 1 && attributes[0].name == "v1_enum" {
			p.attribute(attributes[0], "v1_enum")
			d.Enum = p.enumeration()
		} else {
			d.Interface = p.iface(attributes)
		}
	default:
		p.fail("expected declaration")
	}
	return d
}
func (p *parser) enumeration() *types.EnumDeclaration {
	p.expect("typedef")
	p.expect("enum")
	e := &types.EnumDeclaration{Name: p.ident()}
	p.expect("{")
	for p.err == nil && p.current.text != "}" {
		value := &types.EnumValue{Key: p.ident()}
		if p.take("=") {
			raw := p.current
			if raw.kind != 'n' && raw.kind != 'h' {
				p.fail("expected enum integer")
			}
			p.next()
			value.Value = &types.EnumValueDecl{Value: raw.text}
			if p.take("<") {
				p.expect("<")
				number := p.value('n', "shift count")
				parsed, err := strconv.ParseInt(number, 0, 0)
				n := int(parsed)
				if err != nil {
					p.fail("invalid shift count")
				}
				value.Value.LeftShift = &n
			}
		}
		p.take(",")
		e.Values = append(e.Values, value)
	}
	if len(e.Values) == 0 {
		p.fail("expected enum value")
	}
	p.expect("}")
	p.ident()
	p.expect(";")
	return e
}
func (p *parser) structure() *types.StructDeclaration {
	p.expect("typedef")
	p.expect("struct")
	s := &types.StructDeclaration{Name: p.ident()}
	p.expect("{")
	for p.err == nil && p.current.text != "}" {
		typ := p.ident()
		if typ != "UINT32" && typ != "BOOL" && typ != "BYTE" {
			p.fail("unsupported struct field type")
		}
		s.Fields = append(s.Fields, &types.StructField{Type: typ, Name: p.ident()})
		p.expect(";")
	}
	if len(s.Fields) == 0 {
		p.fail("expected struct field")
	}
	p.expect("}")
	p.ident()
	p.expect(";")
	return s
}
func (p *parser) iface(attributes []attribute) *types.InterfaceDeclaration {
	i := &types.InterfaceDeclaration{Header: &types.InterfaceHeader{}}
	if len(attributes) != 3 {
		p.fail("expected interface attributes")
		return i
	}
	uuid := types.UUID(p.uuid(attributes[0]))
	i.Header.UUID = &uuid
	p.attribute(attributes[1], "object")
	p.attribute(attributes[2], "pointer_default", "unique")
	p.expect("interface")
	i.Name = p.ident()
	p.expect(":")
	i.BaseClass = p.ident()
	p.expect("{")
	for p.err == nil && p.current.text != "}" {
		i.Methods = append(i.Methods, p.method())
	}
	if len(i.Methods) == 0 {
		p.fail("expected method")
	}
	p.expect("}")
	return i
}
func (p *parser) method() *types.InterfaceMethod {
	m := &types.InterfaceMethod{}
	if p.current.text == "[" {
		a := p.attributes()
		if len(a) != 1 || len(a[0].args) != 0 || (a[0].name != "propget" && a[0].name != "propput") {
			p.fail("expected property attribute")
		} else {
			prop := types.Prop("")
			_ = prop.Capture([]string{a[0].name})
			m.Prop = &prop
		}
	}
	m.ReturnType = p.ident()
	_ = m.Name.Capture([]string{p.ident()})
	p.expect("(")
	for p.err == nil && p.current.text != ")" {
		m.Params = append(m.Params, p.param())
	}
	p.expect(")")
	p.expect(";")
	return m
}
func (p *parser) param() *types.Param {
	param := &types.Param{}
	if p.current.text == "[" {
		a := p.attributes()
		direction := &types.Direction{}
		param.Direction = direction
		if len(a) < 1 || len(a) > 2 {
			p.fail("expected parameter direction")
			return param
		}
		if a[0].name != "in" && a[0].name != "out" || len(a[0].args) != 0 {
			p.fail("expected in or out")
		}
		direction.Dir = a[0].name
		if len(a) == 2 {
			switch a[1].name {
			case "retval":
				p.attribute(a[1], "retval")
				direction.Retval = "retval"
			case "size_is":
				if len(a[1].args) != 1 || a[1].args[0].kind != 'i' {
					p.fail("expected size_is identifier")
				} else {
					direction.Retval = "size_is(" + a[1].args[0].text + ")"
				}
			default:
				p.fail("unsupported direction attribute")
			}
		}
	}
	param.Type = p.ident()
	if p.take("const") {
		param.Const = "const"
	}
	for p.take("*") {
		param.Pointer += "*"
	}
	param.Name = p.ident()
	p.take(",")
	return param
}
