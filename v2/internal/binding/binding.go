package binding

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
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
	tsPrefix            string
	tsSuffix            string
	obfuscate           bool
}

// NewBindings returns a new Bindings object
func NewBindings(logger *logger.Logger, structPointersToBind []interface{}, exemptions []interface{}, obfuscate bool) *Bindings {
	result := &Bindings{
		db:                  newDB(),
		logger:              logger.CustomLogger("Bindings"),
		structsToGenerateTS: make(map[string]map[string]interface{}),
		obfuscate:           obfuscate,
	}

	for _, exemption := range exemptions {
		if exemption == nil {
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

func (b *Bindings) GenerateModels() ([]byte, error) {
	models := map[string]string{}
	var seen slicer.StringSlicer
	allStructNames := b.getAllStructNames()
	allStructNames.Sort()
	for packageName, structsToGenerate := range b.structsToGenerateTS {
		thisPackageCode := ""
		w := typescriptify.New()
		w.WithPrefix(b.tsPrefix)
		w.WithSuffix(b.tsSuffix)
		w.Namespace = packageName
		w.WithBackupDir("")
		w.KnownStructs = allStructNames
		// sort the structs
		var structNames []string
		for structName := range structsToGenerate {
			structNames = append(structNames, structName)
		}
		sort.Strings(structNames)
		for _, structName := range structNames {
			fqstructname := packageName + "." + structName
			if seen.Contains(fqstructname) {
				continue
			}
			structInterface := structsToGenerate[structName]
			w.Add(structInterface)
		}
		str, err := w.Convert(nil)
		if err != nil {
			return nil, err
		}
		thisPackageCode += str
		seen.AddSlice(w.GetGeneratedStructs())
		models[packageName] = thisPackageCode
	}

	// Sort the package names first to make the output deterministic
	sortedPackageNames := make([]string, 0)
	for packageName := range models {
		sortedPackageNames = append(sortedPackageNames, packageName)
	}
	sort.Strings(sortedPackageNames)

	var modelsData bytes.Buffer
	for _, packageName := range sortedPackageNames {
		modelData := models[packageName]
		if strings.TrimSpace(modelData) == "" {
			continue
		}
		modelsData.WriteString("export namespace " + packageName + " {\n")
		sc := bufio.NewScanner(strings.NewReader(modelData))
		for sc.Scan() {
			modelsData.WriteString("\t" + sc.Text() + "\n")
		}
		modelsData.WriteString("\n}\n\n")
	}
	return modelsData.Bytes(), nil
}

func (b *Bindings) WriteModels(modelsDir string) error {

	modelsData, err := b.GenerateModels()
	if err != nil {
		return err
	}
	// Don't write if we don't have anything
	if len(modelsData) == 0 {
		return nil
	}

	filename := filepath.Join(modelsDir, "models.ts")
	err = os.WriteFile(filename, modelsData, 0755)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bindings) AddStructToGenerateTS(packageName string, structName string, s interface{}) {
	if b.structsToGenerateTS[packageName] == nil {
		b.structsToGenerateTS[packageName] = make(map[string]interface{})
	}
	if b.structsToGenerateTS[packageName][structName] != nil {
		return
	}
	b.structsToGenerateTS[packageName][structName] = s

	// Iterate this struct and add any struct field references
	structType := reflect.TypeOf(s)
	if hasElements(structType) {
		structType = structType.Elem()
	}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.Anonymous {
			continue
		}
		kind := field.Type.Kind()
		if kind == reflect.Struct {
			if !field.IsExported() {
				continue
			}
			fqname := field.Type.String()
			sNameSplit := strings.Split(fqname, ".")
			if len(sNameSplit) < 2 {
				continue
			}
			sName := sNameSplit[1]
			pName := getPackageName(fqname)
			a := reflect.New(field.Type)
			if b.hasExportedJSONFields(field.Type) {
				s := reflect.Indirect(a).Interface()
				b.AddStructToGenerateTS(pName, sName, s)
			}
		} else if hasElements(field.Type) && field.Type.Elem().Kind() == reflect.Struct {
			if !field.IsExported() {
				continue
			}
			fqname := field.Type.Elem().String()
			sNameSplit := strings.Split(fqname, ".")
			if len(sNameSplit) < 2 {
				continue
			}
			sName := sNameSplit[1]
			pName := getPackageName(fqname)
			typ := field.Type.Elem()
			a := reflect.New(typ)
			if b.hasExportedJSONFields(typ) {
				s := reflect.Indirect(a).Interface()
				b.AddStructToGenerateTS(pName, sName, s)
			}
		}
	}
}

func (b *Bindings) SetTsPrefix(prefix string) *Bindings {
	b.tsPrefix = prefix
	return b
}

func (b *Bindings) SetTsSuffix(postfix string) *Bindings {
	b.tsSuffix = postfix
	return b
}

func (b *Bindings) getAllStructNames() *slicer.StringSlicer {
	var result slicer.StringSlicer
	for packageName, structsToGenerate := range b.structsToGenerateTS {
		for structName := range structsToGenerate {
			result.Add(packageName + "." + structName)
		}
	}
	return &result
}

func (b *Bindings) hasExportedJSONFields(typeOf reflect.Type) bool {
	for i := 0; i < typeOf.NumField(); i++ {
		jsonFieldName := ""
		f := typeOf.Field(i)
		jsonTag := f.Tag.Get("json")
		if len(jsonTag) == 0 {
			continue
		}
		jsonTagParts := strings.Split(jsonTag, ",")
		if len(jsonTagParts) > 0 {
			jsonFieldName = jsonTagParts[0]
		}
		for _, t := range jsonTagParts {
			if t == "-" {
				continue
			}
		}
		if jsonFieldName != "" {
			return true
		}
	}
	return false
}
