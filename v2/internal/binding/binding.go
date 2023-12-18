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
	enumsToGenerateTS   map[string]map[string]interface{}
	tsPrefix            string
	tsSuffix            string
	tsInterface         bool
	obfuscate           bool
}

// NewBindings returns a new Bindings object
func NewBindings(logger *logger.Logger, structPointersToBind []interface{}, exemptions []interface{}, obfuscate bool, enumsToBind []interface{}) *Bindings {
	result := &Bindings{
		db:                  newDB(),
		logger:              logger.CustomLogger("Bindings"),
		structsToGenerateTS: make(map[string]map[string]interface{}),
		enumsToGenerateTS:   make(map[string]map[string]interface{}),
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

	for _, enum := range enumsToBind {
		result.AddEnumToGenerateTS(enum)
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
	var seenEnumsPackages slicer.StringSlicer
	allStructNames := b.getAllStructNames()
	allStructNames.Sort()
	allEnumNames := b.getAllEnumNames()
	allEnumNames.Sort()
	for packageName, structsToGenerate := range b.structsToGenerateTS {
		thisPackageCode := ""
		w := typescriptify.New()
		w.WithPrefix(b.tsPrefix)
		w.WithSuffix(b.tsSuffix)
		w.WithInterface(b.tsInterface)
		w.Namespace = packageName
		w.WithBackupDir("")
		w.KnownStructs = allStructNames
		w.KnownEnums = allEnumNames
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

		// if we have enums for this package, add them as well
		var enums, enumsExist = b.enumsToGenerateTS[packageName]
		if enumsExist {
			for enumName, enum := range enums {
				fqemumname := packageName + "." + enumName
				if seen.Contains(fqemumname) {
					continue
				}
				w.AddEnum(enum)
			}
			seenEnumsPackages.Add(packageName)
		}

		str, err := w.Convert(nil)
		if err != nil {
			return nil, err
		}
		thisPackageCode += str
		seen.AddSlice(w.GetGeneratedStructs())
		models[packageName] = thisPackageCode
	}

	// Add outstanding enums to the models that were not in packages with structs
	for packageName, enumsToGenerate := range b.enumsToGenerateTS {
		if seenEnumsPackages.Contains(packageName) {
			continue
		}

		thisPackageCode := ""
		w := typescriptify.New()
		w.WithPrefix(b.tsPrefix)
		w.WithSuffix(b.tsSuffix)
		w.WithInterface(b.tsInterface)
		w.Namespace = packageName
		w.WithBackupDir("")

		for enumName, enum := range enumsToGenerate {
			fqemumname := packageName + "." + enumName
			if seen.Contains(fqemumname) {
				continue
			}
			w.AddEnum(enum)
		}
		str, err := w.Convert(nil)
		if err != nil {
			return nil, err
		}
		thisPackageCode += str
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
	err = os.WriteFile(filename, modelsData, 0o755)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bindings) AddEnumToGenerateTS(e interface{}) {
	enumType := reflect.TypeOf(e)

	var packageName string
	var enumName string
	// enums should be represented as array of all possible values
	if hasElements(enumType) {
		enum := enumType.Elem()
		// simple enum represented by struct with Value/TSName fields
		if enum.Kind() == reflect.Struct {
			_, tsNamePresented := enum.FieldByName("TSName")
			enumT, valuePresented := enum.FieldByName("Value")
			if tsNamePresented && valuePresented {
				packageName = getPackageName(enumT.Type.String())
				enumName = enumT.Type.Name()
			} else {
				return
			}
			// otherwise expecting implementation with TSName() https://github.com/tkrajina/typescriptify-golang-structs#enums-with-tsname
		} else {
			packageName = getPackageName(enumType.Elem().String())
			enumName = enumType.Elem().Name()
		}
		if b.enumsToGenerateTS[packageName] == nil {
			b.enumsToGenerateTS[packageName] = make(map[string]interface{})
		}
		if b.enumsToGenerateTS[packageName][enumName] != nil {
			return
		}
		b.enumsToGenerateTS[packageName][enumName] = e
	}
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

func (b *Bindings) SetOutputType(outputType string) *Bindings {
	if outputType == "interfaces" {
		b.tsInterface = true
	}
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

func (b *Bindings) getAllEnumNames() *slicer.StringSlicer {
	var result slicer.StringSlicer
	for packageName, enumsToGenerate := range b.enumsToGenerateTS {
		for enumName := range enumsToGenerate {
			result.Add(packageName + "." + enumName)
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
