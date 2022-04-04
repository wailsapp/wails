package binding

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v2/internal/typescriptify"

	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/internal/logger"
)

type Bindings struct {
	db         *DB
	logger     logger.CustomLogger
	exemptions slicer.StringSlicer

	structsToGenerateTS map[string]map[string]interface{}
}

// NewBindings returns a new Bindings object
func NewBindings(logger *logger.Logger, structPointersToBind []interface{}, exemptions []interface{}) *Bindings {
	result := &Bindings{
		db:                  newDB(),
		logger:              logger.CustomLogger("Bindings"),
		structsToGenerateTS: make(map[string]map[string]interface{}),
	}

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

func (b *Bindings) DB() *DB {
	return b.db
}

func (b *Bindings) ToJSON() (string, error) {
	return b.db.ToJSON()
}

func (b *Bindings) WriteModels(modelsDir string) error {
	models := map[string]string{}
	for packageName, structsToGenerate := range b.structsToGenerateTS {
		thisPackageCode := ""
		for _, structInterface := range structsToGenerate {
			w := typescriptify.New()
			w.WithBackupDir("")
			w.Add(structInterface)
			str, err := w.Convert(nil)
			if err != nil {
				return err
			}
			thisPackageCode += str
		}
		models[packageName] = thisPackageCode
	}

	var modelsData bytes.Buffer
	for packageName, modelData := range models {
		modelsData.WriteString("export namespace " + packageName + " {\n")
		sc := bufio.NewScanner(strings.NewReader(modelData))
		for sc.Scan() {
			modelsData.WriteString("\t" + sc.Text() + "\n")
		}
		modelsData.WriteString("\n}\n\n")
	}

	filename := filepath.Join(modelsDir, "models.ts")
	err := os.WriteFile(filename, modelsData.Bytes(), 0755)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bindings) AddStructToGenerateTS(packageName string, structName string, s interface{}) {
	println("Adding struct:", packageName, structName)
	if b.structsToGenerateTS[packageName] == nil {
		b.structsToGenerateTS[packageName] = make(map[string]interface{})
	}
	b.structsToGenerateTS[packageName][structName] = s
}
