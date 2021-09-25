package binding

import (
	"fmt"
	"github.com/leaanthony/typescriptify-golang-structs/typescriptify"
	"reflect"
	"runtime"
	"strings"

	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/internal/logger"
)

type Bindings struct {
	db         *DB
	logger     logger.CustomLogger
	exemptions slicer.StringSlicer

	// Typescript writer
	converter *typescriptify.TypeScriptify
}

// NewBindings returns a new Bindings object
func NewBindings(logger *logger.Logger, structPointersToBind []interface{}, exemptions []interface{}) *Bindings {
	result := &Bindings{
		db:        newDB(),
		logger:    logger.CustomLogger("Bindings"),
		converter: typescriptify.New(),
	}

	// No backups
	result.converter.WithBackupDir("")

	// Hack for TS compilation error
	result.converter.AddImport("export {};")

	for _, exemption := range exemptions {
		if exemptions == nil {
			continue
		}
		name := runtime.FuncForPC(reflect.ValueOf(exemption).Pointer()).Name()
		// Yuk yuk yuk! Is there a better way?
		name = strings.TrimSuffix(name, "-fm")
		result.exemptions.Add(name)
	}

	// Add the structs to bind
	for _, ptr := range structPointersToBind {
		err := result.Add(ptr)
		if err != nil {
			logger.Fatal("Error during binding: " + err.Error())
		}
	}

	return result
}

// Add the given struct methods to the Bindings
func (b *Bindings) Add(structPtr interface{}) error {

	methods, err := b.getMethods(structPtr)
	if err != nil {
		return fmt.Errorf("cannot bind value to app: %s", err.Error())
	}

	for _, method := range methods {
		splitName := strings.Split(method.Name, ".")
		packageName := splitName[0]
		structName := splitName[1]
		methodName := splitName[2]

		// Add it as a regular method
		b.db.AddMethod(packageName, structName, methodName, method)
	}
	return nil
}

func (b *Bindings) WriteTS(filename string) error {
	return b.converter.ConvertToFile(filename)
}

func (b *Bindings) DB() *DB {
	return b.db
}

func (b *Bindings) ToJSON() (string, error) {
	return b.db.ToJSON()
}
