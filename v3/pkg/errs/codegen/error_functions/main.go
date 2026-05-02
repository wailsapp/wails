package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"text/template"
)

const tmpl = `package errs

import (
	"fmt"
)

type wailsError struct {
	cause           error
	msg             string
	errorType       ErrorType
}

func (w *wailsError) Cause() error          { return w.cause }
func (w *wailsError) Error() string {
	errMsg := fmt.Sprintf("%s: %s", w.errorType, w.msg)
	if w.cause != nil {
		return fmt.Sprintf("%s: %s", errMsg, w.cause.Error())
	}
	return errMsg
}
func (w *wailsError) Msg() string           { return w.msg }
func (w *wailsError) ErrorType() ErrorType  { return w.errorType }
func (w *wailsError) Unwrap() error  { return w.cause }

{{ range . }}
func New{{ . }}f(message string, args ...any) error {
	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:           nil,
		msg:             msg,
		errorType:       {{ . }},
	}
}

func Wrap{{ . }}f(err error, message string, args ...any) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:           err,
		msg:             msg,
		errorType:       {{ . }},
	}
}

func Is{{ . }}(err error) bool {
	return Is(err, {{ . }})
}

func Has{{ . }}(err error) bool {
	return Has(err, {{ . }})
}
{{ end }}
`

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "errors.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	var errorTypes []string
	ast.Inspect(f, func(n ast.Node) bool {
		v, ok := n.(*ast.ValueSpec)
		if !ok {
			return true
		}
		for _, name := range v.Names {
			errorTypes = append(errorTypes, name.Name)
		}
		return true
	})

	file, err := os.Create("error_functions.gen.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	t := template.Must(template.New("tmpl").Parse(tmpl))
	if err := t.Execute(file, errorTypes); err != nil {
		panic(err)
	}
}
