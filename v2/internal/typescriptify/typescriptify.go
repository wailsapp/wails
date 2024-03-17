package typescriptify

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/leaanthony/slicer"

	"github.com/tkrajina/go-reflector/reflector"
)

const (
	tsTransformTag      = "ts_transform"
	tsType              = "ts_type"
	tsConvertValuesFunc = `convertValues(a: any, classs: any, asMap: boolean = false): any {
	if (!a) {
		return a;
	}
	if (a.slice) {
		return (a as any[]).map(elem => this.convertValues(elem, classs));
	} else if ("object" === typeof a) {
		if (asMap) {
			for (const key of Object.keys(a)) {
				a[key] = new classs(a[key]);
			}
			return a;
		}
		return new classs(a);
	}
	return a;
}`
	jsVariableNameRegex = `^([A-Z]|[a-z]|\$|_)([A-Z]|[a-z]|[0-9]|\$|_)*$`
)

// TypeOptions overrides options set by `ts_*` tags.
type TypeOptions struct {
	TSType      string
	TSTransform string
}

// StructType stores settings for transforming one Golang struct.
type StructType struct {
	Type         reflect.Type
	FieldOptions map[reflect.Type]TypeOptions
}

func NewStruct(i interface{}) *StructType {
	return &StructType{
		Type: reflect.TypeOf(i),
	}
}

func (st *StructType) WithFieldOpts(i interface{}, opts TypeOptions) *StructType {
	if st.FieldOptions == nil {
		st.FieldOptions = map[reflect.Type]TypeOptions{}
	}
	var typ reflect.Type
	if ty, is := i.(reflect.Type); is {
		typ = ty
	} else {
		typ = reflect.TypeOf(i)
	}
	st.FieldOptions[typ] = opts
	return st
}

type EnumType struct {
	Type reflect.Type
}

type enumElement struct {
	value interface{}
	name  string
}

type TypeScriptify struct {
	Prefix            string
	Suffix            string
	Indent            string
	CreateFromMethod  bool
	CreateConstructor bool
	BackupDir         string // If empty no backup
	DontExport        bool
	CreateInterface   bool
	customImports     []string

	structTypes []StructType
	enumTypes   []EnumType
	enums       map[reflect.Type][]enumElement
	kinds       map[reflect.Kind]string

	fieldTypeOptions map[reflect.Type]TypeOptions

	// throwaway, used when converting
	alreadyConverted map[string]bool

	Namespace    string
	KnownStructs *slicer.StringSlicer
	KnownEnums   *slicer.StringSlicer
}

func New() *TypeScriptify {
	result := new(TypeScriptify)
	result.Indent = "\t"
	result.BackupDir = "."

	kinds := make(map[reflect.Kind]string)

	kinds[reflect.Bool] = "boolean"
	kinds[reflect.Interface] = "any"

	kinds[reflect.Int] = "number"
	kinds[reflect.Int8] = "number"
	kinds[reflect.Int16] = "number"
	kinds[reflect.Int32] = "number"
	kinds[reflect.Int64] = "number"
	kinds[reflect.Uint] = "number"
	kinds[reflect.Uint8] = "number"
	kinds[reflect.Uint16] = "number"
	kinds[reflect.Uint32] = "number"
	kinds[reflect.Uint64] = "number"
	kinds[reflect.Float32] = "number"
	kinds[reflect.Float64] = "number"

	kinds[reflect.String] = "string"

	result.kinds = kinds

	result.Indent = "    "
	result.CreateFromMethod = true
	result.CreateConstructor = true

	return result
}

func (t *TypeScriptify) deepFields(typeOf reflect.Type) []reflect.StructField {
	fields := make([]reflect.StructField, 0)

	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}

	if typeOf.Kind() != reflect.Struct {
		return fields
	}

	for i := 0; i < typeOf.NumField(); i++ {
		f := typeOf.Field(i)
		kind := f.Type.Kind()
		isPointer := kind == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct
		if f.Anonymous && kind == reflect.Struct {
			// fmt.Println(v.Interface())
			fields = append(fields, t.deepFields(f.Type)...)
		} else if f.Anonymous && isPointer {
			// fmt.Println(v.Interface())
			fields = append(fields, t.deepFields(f.Type.Elem())...)
		} else {
			// Check we have a json tag
			jsonTag := t.getJSONFieldName(f, isPointer)
			if jsonTag != "" {
				fields = append(fields, f)
			}
		}
	}

	return fields
}

func (ts TypeScriptify) logf(depth int, s string, args ...interface{}) {
	fmt.Printf(strings.Repeat("   ", depth)+s+"\n", args...)
}

// ManageType can define custom options for fields of a specified type.
//
// This can be used instead of setting ts_type and ts_transform for all fields of a certain type.
func (t *TypeScriptify) ManageType(fld interface{}, opts TypeOptions) *TypeScriptify {
	var typ reflect.Type
	switch t := fld.(type) {
	case reflect.Type:
		typ = t
	default:
		typ = reflect.TypeOf(fld)
	}
	if t.fieldTypeOptions == nil {
		t.fieldTypeOptions = map[reflect.Type]TypeOptions{}
	}
	t.fieldTypeOptions[typ] = opts
	return t
}

func (t *TypeScriptify) GetGeneratedStructs() []string {
	var result []string
	for key := range t.alreadyConverted {
		result = append(result, key)
	}
	return result
}

func (t *TypeScriptify) WithCreateFromMethod(b bool) *TypeScriptify {
	t.CreateFromMethod = b
	return t
}

func (t *TypeScriptify) WithInterface(b bool) *TypeScriptify {
	t.CreateInterface = b
	return t
}

func (t *TypeScriptify) WithConstructor(b bool) *TypeScriptify {
	t.CreateConstructor = b
	return t
}

func (t *TypeScriptify) WithIndent(i string) *TypeScriptify {
	t.Indent = i
	return t
}

func (t *TypeScriptify) WithBackupDir(b string) *TypeScriptify {
	t.BackupDir = b
	return t
}

func (t *TypeScriptify) WithPrefix(p string) *TypeScriptify {
	t.Prefix = p
	return t
}

func (t *TypeScriptify) WithSuffix(s string) *TypeScriptify {
	t.Suffix = s
	return t
}

func (t *TypeScriptify) Add(obj interface{}) *TypeScriptify {
	switch ty := obj.(type) {
	case StructType:
		t.structTypes = append(t.structTypes, ty)
	case *StructType:
		t.structTypes = append(t.structTypes, *ty)
	case reflect.Type:
		t.AddType(ty)
	default:
		t.AddType(reflect.TypeOf(obj))
	}
	return t
}

func (t *TypeScriptify) AddType(typeOf reflect.Type) *TypeScriptify {
	t.structTypes = append(t.structTypes, StructType{Type: typeOf})
	return t
}

func (t *typeScriptClassBuilder) AddMapField(fieldName string, field reflect.StructField) {
	keyType := field.Type.Key()
	valueType := field.Type.Elem()
	valueTypeName := valueType.Name()
	if name, ok := t.types[valueType.Kind()]; ok {
		valueTypeName = name
	}
	if valueType.Kind() == reflect.Array || valueType.Kind() == reflect.Slice {
		valueTypeName = valueType.Elem().Name() + "[]"
	}
	if valueType.Kind() == reflect.Ptr {
		valueTypeName = valueType.Elem().Name()
	}
	if valueType.Kind() == reflect.Struct && differentNamespaces(t.namespace, valueType) {
		valueTypeName = valueType.String()
	}
	strippedFieldName := strings.ReplaceAll(fieldName, "?", "")
	isOptional := strings.HasSuffix(fieldName, "?")

	keyTypeStr := ""
	// Key should always be a JS primitive. JS will read it as a string either way.
	if typeStr, isSimple := t.types[keyType.Kind()]; isSimple {
		keyTypeStr = typeStr
	} else {
		keyTypeStr = t.types[reflect.String]
	}

	var dotField string
	if regexp.MustCompile(jsVariableNameRegex).Match([]byte(strippedFieldName)) {
		dotField = fmt.Sprintf(".%s", strippedFieldName)
	} else {
		dotField = fmt.Sprintf(`["%s"]`, strippedFieldName)
		if isOptional {
			fieldName = fmt.Sprintf(`"%s"?`, strippedFieldName)
		}
	}
	t.fields = append(t.fields, fmt.Sprintf("%s%s: {[key: %s]: %s};", t.indent, fieldName, keyTypeStr, valueTypeName))
	if valueType.Kind() == reflect.Struct {
		t.constructorBody = append(t.constructorBody, fmt.Sprintf("%s%sthis%s = this.convertValues(source[\"%s\"], %s, true);", t.indent, t.indent, dotField, strippedFieldName, t.prefix+valueTypeName+t.suffix))
	} else {
		t.constructorBody = append(t.constructorBody, fmt.Sprintf("%s%sthis%s = source[\"%s\"];", t.indent, t.indent, dotField, strippedFieldName))
	}
}

func (t *TypeScriptify) AddEnum(values interface{}) *TypeScriptify {
	if t.enums == nil {
		t.enums = map[reflect.Type][]enumElement{}
	}
	items := reflect.ValueOf(values)
	if items.Kind() != reflect.Slice {
		panic(fmt.Sprintf("Values for %T isn't a slice", values))
	}

	var elements []enumElement
	for i := 0; i < items.Len(); i++ {
		item := items.Index(i)

		var el enumElement
		if item.Kind() == reflect.Struct {
			r := reflector.New(item.Interface())
			val, err := r.Field("Value").Get()
			if err != nil {
				panic(fmt.Sprint("missing Type field in ", item.Type().String()))
			}
			name, err := r.Field("TSName").Get()
			if err != nil {
				panic(fmt.Sprint("missing TSName field in ", item.Type().String()))
			}
			el.value = val
			el.name = name.(string)
		} else {
			el.value = item.Interface()
			if tsNamer, is := item.Interface().(TSNamer); is {
				el.name = tsNamer.TSName()
			} else {
				panic(fmt.Sprint(item.Type().String(), " has no TSName method"))
			}
		}

		elements = append(elements, el)
	}
	ty := reflect.TypeOf(elements[0].value)
	t.enums[ty] = elements
	t.enumTypes = append(t.enumTypes, EnumType{Type: ty})

	return t
}

// AddEnumValues is deprecated, use `AddEnum()`
func (t *TypeScriptify) AddEnumValues(typeOf reflect.Type, values interface{}) *TypeScriptify {
	t.AddEnum(values)
	return t
}

func (t *TypeScriptify) Convert(customCode map[string]string) (string, error) {
	t.alreadyConverted = make(map[string]bool)
	depth := 0

	result := ""
	if len(t.customImports) > 0 {
		// Put the custom imports, i.e.: `import Decimal from 'decimal.js'`
		for _, cimport := range t.customImports {
			result += cimport + "\n"
		}
	}

	for _, enumTyp := range t.enumTypes {
		elements := t.enums[enumTyp.Type]
		typeScriptCode, err := t.convertEnum(depth, enumTyp.Type, elements)
		if err != nil {
			return "", err
		}
		result += "\n" + strings.Trim(typeScriptCode, " "+t.Indent+"\r\n")
	}

	for _, strctTyp := range t.structTypes {
		typeScriptCode, err := t.convertType(depth, strctTyp.Type, customCode)
		if err != nil {
			return "", err
		}
		result += "\n" + strings.Trim(typeScriptCode, " "+t.Indent+"\r\n")
	}
	return result, nil
}

func loadCustomCode(fileName string) (map[string]string, error) {
	result := make(map[string]string)
	f, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return result, nil
		}
		return result, err
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return result, err
	}

	var currentName string
	var currentValue string
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "//[") && strings.HasSuffix(trimmedLine, ":]") {
			currentName = strings.Replace(strings.Replace(trimmedLine, "//[", "", -1), ":]", "", -1)
			currentValue = ""
		} else if trimmedLine == "//[end]" {
			result[currentName] = strings.TrimRight(currentValue, " \t\r\n")
			currentName = ""
			currentValue = ""
		} else if len(currentName) > 0 {
			currentValue += line + "\n"
		}
	}

	return result, nil
}

func (t TypeScriptify) backup(fileName string) error {
	fileIn, err := os.Open(fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// No neet to backup, just return:
		return nil
	}
	defer fileIn.Close()

	bytes, err := ioutil.ReadAll(fileIn)
	if err != nil {
		return err
	}

	_, backupFn := path.Split(fmt.Sprintf("%s-%s.backup", fileName, time.Now().Format("2006-01-02T15_04_05.99")))
	if t.BackupDir != "" {
		backupFn = path.Join(t.BackupDir, backupFn)
	}

	return ioutil.WriteFile(backupFn, bytes, os.FileMode(0o700))
}

func (t TypeScriptify) ConvertToFile(fileName string, packageName string) error {
	if len(t.BackupDir) > 0 {
		err := t.backup(fileName)
		if err != nil {
			return err
		}
	}

	customCode, err := loadCustomCode(fileName)
	if err != nil {
		return err
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	converted, err := t.Convert(customCode)
	if err != nil {
		return err
	}

	var lines []string
	sc := bufio.NewScanner(strings.NewReader(converted))
	for sc.Scan() {
		lines = append(lines, "\t"+sc.Text())
	}

	converted = "export namespace " + packageName + " {\n"
	converted += strings.Join(lines, "\n")
	converted += "\n}\n"

	if _, err := f.WriteString("/* Do not change, this code is generated from Golang structs */\n\n"); err != nil {
		return err
	}
	if _, err := f.WriteString(converted); err != nil {
		return err
	}
	if err != nil {
		return err
	}

	return nil
}

type TSNamer interface {
	TSName() string
}

func (t *TypeScriptify) convertEnum(depth int, typeOf reflect.Type, elements []enumElement) (string, error) {
	t.logf(depth, "Converting enum %s", typeOf.String())
	if _, found := t.alreadyConverted[typeOf.String()]; found { // Already converted
		return "", nil
	}
	t.alreadyConverted[typeOf.String()] = true

	entityName := t.Prefix + typeOf.Name() + t.Suffix
	result := "enum " + entityName + " {\n"

	for _, val := range elements {
		result += fmt.Sprintf("%s%s = %#v,\n", t.Indent, val.name, val.value)
	}

	result += "}"

	if !t.DontExport {
		result = "export " + result
	}

	return result, nil
}

func (t *TypeScriptify) getFieldOptions(structType reflect.Type, field reflect.StructField) TypeOptions {
	// By default use options defined by tags:
	opts := TypeOptions{TSTransform: field.Tag.Get(tsTransformTag), TSType: field.Tag.Get(tsType)}

	overrides := []TypeOptions{}

	// But there is maybe an struct-specific override:
	for _, strct := range t.structTypes {
		if strct.FieldOptions == nil {
			continue
		}
		if strct.Type == structType {
			if fldOpts, found := strct.FieldOptions[field.Type]; found {
				overrides = append(overrides, fldOpts)
			}
		}
	}

	if fldOpts, found := t.fieldTypeOptions[field.Type]; found {
		overrides = append(overrides, fldOpts)
	}

	for _, o := range overrides {
		if o.TSTransform != "" {
			opts.TSTransform = o.TSTransform
		}
		if o.TSType != "" {
			opts.TSType = o.TSType
		}
	}

	return opts
}

func (t *TypeScriptify) getJSONFieldName(field reflect.StructField, isPtr bool) string {
	jsonFieldName := ""
	jsonTag := field.Tag.Get("json")
	if len(jsonTag) > 0 {
		jsonTagParts := strings.Split(jsonTag, ",")
		if len(jsonTagParts) > 0 {
			jsonFieldName = strings.Trim(jsonTagParts[0], t.Indent)
		}
		hasOmitEmpty := false
		ignored := false
		for _, t := range jsonTagParts {
			if t == "" {
				break
			}
			if t == "omitempty" {
				hasOmitEmpty = true
				break
			}
			if t == "-" {
				ignored = true
				break
			}
		}
		if !ignored && isPtr || hasOmitEmpty {
			jsonFieldName = fmt.Sprintf("%s?", jsonFieldName)
		}
	}
	return jsonFieldName
}

func (t *TypeScriptify) convertType(depth int, typeOf reflect.Type, customCode map[string]string) (string, error) {
	if _, found := t.alreadyConverted[typeOf.String()]; found { // Already converted
		return "", nil
	}
	fields := t.deepFields(typeOf)
	if len(fields) == 0 {
		return "", nil
	}
	t.logf(depth, "Converting type %s", typeOf.String())
	if differentNamespaces(t.Namespace, typeOf) {
		return "", nil
	}

	t.alreadyConverted[typeOf.String()] = true

	entityName := t.Prefix + typeOf.Name() + t.Suffix

	if typeClashWithReservedKeyword(entityName) {
		warnAboutTypesClash(entityName)
	}

	result := ""
	if t.CreateInterface {
		result += fmt.Sprintf("interface %s {\n", entityName)
	} else {
		result += fmt.Sprintf("class %s {\n", entityName)
	}
	if !t.DontExport {
		result = "export " + result
	}
	builder := typeScriptClassBuilder{
		types:     t.kinds,
		indent:    t.Indent,
		prefix:    t.Prefix,
		suffix:    t.Suffix,
		namespace: t.Namespace,
	}

	for _, field := range fields {
		isPtr := field.Type.Kind() == reflect.Ptr
		if isPtr {
			field.Type = field.Type.Elem()
		}
		jsonFieldName := t.getJSONFieldName(field, isPtr)
		if len(jsonFieldName) == 0 || jsonFieldName == "-" {
			continue
		}

		var err error
		fldOpts := t.getFieldOptions(typeOf, field)
		if fldOpts.TSTransform != "" {
			t.logf(depth, "- simple field %s.%s", typeOf.Name(), field.Name)
			err = builder.AddSimpleField(jsonFieldName, field, fldOpts)
		} else if _, isEnum := t.enums[field.Type]; isEnum {
			t.logf(depth, "- enum field %s.%s", typeOf.Name(), field.Name)
			builder.AddEnumField(jsonFieldName, field)
		} else if fldOpts.TSType != "" { // Struct:
			t.logf(depth, "- simple field %s.%s", typeOf.Name(), field.Name)
			err = builder.AddSimpleField(jsonFieldName, field, fldOpts)
		} else if field.Type.Kind() == reflect.Struct { // Struct:
			t.logf(depth, "- struct %s.%s (%s)", typeOf.Name(), field.Name, field.Type.String())

			// Anonymous structures is ignored
			// It is possible to generate them but hard to generate correct name
			if field.Type.Name() != "" {
				typeScriptChunk, err := t.convertType(depth+1, field.Type, customCode)
				if err != nil {
					return "", err
				}
				if typeScriptChunk != "" {
					result = typeScriptChunk + "\n" + result
				}
			}

			isKnownType := t.KnownStructs.Contains(getStructFQN(field.Type.String()))
			println("KnownStructs:", t.KnownStructs.Join("\t"))
			println(getStructFQN(field.Type.String()))
			builder.AddStructField(jsonFieldName, field, !isKnownType)
		} else if field.Type.Kind() == reflect.Map {
			t.logf(depth, "- map field %s.%s", typeOf.Name(), field.Name)
			// Also convert map key types if needed
			var keyTypeToConvert reflect.Type
			switch field.Type.Key().Kind() {
			case reflect.Struct:
				keyTypeToConvert = field.Type.Key()
			case reflect.Ptr:
				keyTypeToConvert = field.Type.Key().Elem()
			}
			if keyTypeToConvert != nil {
				typeScriptChunk, err := t.convertType(depth+1, keyTypeToConvert, customCode)
				if err != nil {
					return "", err
				}
				if typeScriptChunk != "" {
					result = typeScriptChunk + "\n" + result
				}
			}
			// Also convert map value types if needed
			var valueTypeToConvert reflect.Type
			switch field.Type.Elem().Kind() {
			case reflect.Struct:
				valueTypeToConvert = field.Type.Elem()
			case reflect.Ptr:
				valueTypeToConvert = field.Type.Elem().Elem()
			}
			if valueTypeToConvert != nil {
				typeScriptChunk, err := t.convertType(depth+1, valueTypeToConvert, customCode)
				if err != nil {
					return "", err
				}
				if typeScriptChunk != "" {
					result = typeScriptChunk + "\n" + result
				}
			}

			builder.AddMapField(jsonFieldName, field)
		} else if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array { // Slice:
			if field.Type.Elem().Kind() == reflect.Ptr { // extract ptr type
				field.Type = field.Type.Elem()
			}

			arrayDepth := 1
			for field.Type.Elem().Kind() == reflect.Slice { // Slice of slices:
				field.Type = field.Type.Elem()
				arrayDepth++
			}

			if field.Type.Elem().Kind() == reflect.Struct { // Slice of structs:
				t.logf(depth, "- struct slice %s.%s (%s)", typeOf.Name(), field.Name, field.Type.String())
				typeScriptChunk, err := t.convertType(depth+1, field.Type.Elem(), customCode)
				if err != nil {
					return "", err
				}
				if typeScriptChunk != "" {
					result = typeScriptChunk + "\n" + result
				}
				builder.AddArrayOfStructsField(jsonFieldName, field, arrayDepth)
			} else { // Slice of simple fields:
				t.logf(depth, "- slice field %s.%s", typeOf.Name(), field.Name)
				err = builder.AddSimpleArrayField(jsonFieldName, field, arrayDepth, fldOpts)
			}
		} else { // Simple field:
			t.logf(depth, "- simple field %s.%s", typeOf.Name(), field.Name)
			// check if type is in known enum. If so, then replace TStype with enum name to avoid missing types
			isKnownEnum := t.KnownEnums.Contains(getStructFQN(field.Type.String()))
			if isKnownEnum {
				err = builder.AddSimpleField(jsonFieldName, field, TypeOptions{
					TSType:      getStructFQN(field.Type.String()),
					TSTransform: fldOpts.TSTransform,
				})
			} else {
				err = builder.AddSimpleField(jsonFieldName, field, fldOpts)
			}
		}
		if err != nil {
			return "", err
		}
	}

	if t.CreateFromMethod {
		t.CreateConstructor = true
	}

	result += strings.Join(builder.fields, "\n") + "\n"
	if !t.CreateInterface {
		constructorBody := strings.Join(builder.constructorBody, "\n")
		needsConvertValue := strings.Contains(constructorBody, "this.convertValues")
		if t.CreateFromMethod {
			result += fmt.Sprintf("\n%sstatic createFrom(source: any = {}) {\n", t.Indent)
			result += fmt.Sprintf("%s%sreturn new %s(source);\n", t.Indent, t.Indent, entityName)
			result += fmt.Sprintf("%s}\n", t.Indent)
		}
		if t.CreateConstructor {
			result += fmt.Sprintf("\n%sconstructor(source: any = {}) {\n", t.Indent)
			result += t.Indent + t.Indent + "if ('string' === typeof source) source = JSON.parse(source);\n"
			result += constructorBody + "\n"
			result += fmt.Sprintf("%s}\n", t.Indent)
		}
		if needsConvertValue && (t.CreateConstructor || t.CreateFromMethod) {
			result += "\n" + indentLines(strings.ReplaceAll(tsConvertValuesFunc, "\t", t.Indent), 1) + "\n"
		}
	}

	if customCode != nil {
		code := customCode[entityName]
		if len(code) != 0 {
			result += t.Indent + "//[" + entityName + ":]\n" + code + "\n\n" + t.Indent + "//[end]\n"
		}
	}

	result += "}"

	return result, nil
}

func (t *TypeScriptify) AddImport(i string) {
	for _, cimport := range t.customImports {
		if cimport == i {
			return
		}
	}

	t.customImports = append(t.customImports, i)
}

type typeScriptClassBuilder struct {
	types                map[reflect.Kind]string
	indent               string
	fields               []string
	createFromMethodBody []string
	constructorBody      []string
	prefix, suffix       string
	namespace            string
}

func (t *typeScriptClassBuilder) AddSimpleArrayField(fieldName string, field reflect.StructField, arrayDepth int, opts TypeOptions) error {
	fieldType, kind := field.Type.Elem().Name(), field.Type.Elem().Kind()
	typeScriptType := t.types[kind]

	if len(fieldName) > 0 {
		strippedFieldName := strings.ReplaceAll(fieldName, "?", "")
		if len(opts.TSType) > 0 {
			t.addField(fieldName, opts.TSType, false)
			t.addInitializerFieldLine(strippedFieldName, fmt.Sprintf("source[\"%s\"]", strippedFieldName))
			return nil
		} else if len(typeScriptType) > 0 {
			t.addField(fieldName, fmt.Sprint(typeScriptType, strings.Repeat("[]", arrayDepth)), false)
			t.addInitializerFieldLine(strippedFieldName, fmt.Sprintf("source[\"%s\"]", strippedFieldName))
			return nil
		}
	}

	return fmt.Errorf("cannot find type for %s (%s/%s)", kind.String(), fieldName, fieldType)
}

func (t *typeScriptClassBuilder) AddSimpleField(fieldName string, field reflect.StructField, opts TypeOptions) error {
	fieldType, kind := field.Type.Name(), field.Type.Kind()

	typeScriptType := t.types[kind]
	if len(opts.TSType) > 0 {
		typeScriptType = opts.TSType
	}

	if len(typeScriptType) > 0 && len(fieldName) > 0 {
		strippedFieldName := strings.ReplaceAll(fieldName, "?", "")
		t.addField(fieldName, typeScriptType, false)
		if opts.TSTransform == "" {
			t.addInitializerFieldLine(strippedFieldName, fmt.Sprintf("source[\"%s\"]", strippedFieldName))
		} else {
			val := fmt.Sprintf(`source["%s"]`, strippedFieldName)
			expression := strings.Replace(opts.TSTransform, "__VALUE__", val, -1)
			t.addInitializerFieldLine(strippedFieldName, expression)
		}
		return nil
	}

	return fmt.Errorf("cannot find type for %s (%s/%s)", kind.String(), fieldName, fieldType)
}

func (t *typeScriptClassBuilder) AddEnumField(fieldName string, field reflect.StructField) {
	fieldType := field.Type.Name()
	t.addField(fieldName, t.prefix+fieldType+t.suffix, false)
	strippedFieldName := strings.ReplaceAll(fieldName, "?", "")
	t.addInitializerFieldLine(strippedFieldName, fmt.Sprintf("source[\"%s\"]", strippedFieldName))
}

func (t *typeScriptClassBuilder) AddStructField(fieldName string, field reflect.StructField, isAnyType bool) {
	strippedFieldName := strings.ReplaceAll(fieldName, "?", "")
	classname := "null"
	namespace := strings.Split(field.Type.String(), ".")[0]
	fqname := t.prefix + field.Type.Name() + t.suffix
	if namespace != t.namespace {
		fqname = namespace + "." + fqname
	}

	if !isAnyType {
		classname = fqname
	}

	// Anonymous struct
	if field.Type.Name() == "" {
		classname = "Object"
	}

	t.addField(fieldName, fqname, isAnyType)
	t.addInitializerFieldLine(strippedFieldName, fmt.Sprintf("this.convertValues(source[\"%s\"], %s)", strippedFieldName, classname))
}

func (t *typeScriptClassBuilder) AddArrayOfStructsField(fieldName string, field reflect.StructField, arrayDepth int) {
	fieldType := field.Type.Elem().Name()
	if differentNamespaces(t.namespace, field.Type.Elem()) {
		fieldType = field.Type.Elem().String()
	}
	strippedFieldName := strings.ReplaceAll(fieldName, "?", "")
	t.addField(fieldName, fmt.Sprint(t.prefix+fieldType+t.suffix, strings.Repeat("[]", arrayDepth)), false)
	t.addInitializerFieldLine(strippedFieldName, fmt.Sprintf("this.convertValues(source[\"%s\"], %s)", strippedFieldName, t.prefix+fieldType+t.suffix))
}

func (t *typeScriptClassBuilder) addInitializerFieldLine(fld, initializer string) {
	var dotField string
	if regexp.MustCompile(jsVariableNameRegex).Match([]byte(fld)) {
		dotField = fmt.Sprintf(".%s", fld)
	} else {
		dotField = fmt.Sprintf(`["%s"]`, fld)
	}
	t.createFromMethodBody = append(t.createFromMethodBody, fmt.Sprint(t.indent, t.indent, "result", dotField, " = ", initializer, ";"))
	t.constructorBody = append(t.constructorBody, fmt.Sprint(t.indent, t.indent, "this", dotField, " = ", initializer, ";"))
}

func (t *typeScriptClassBuilder) addField(fld, fldType string, isAnyType bool) {
	isOptional := strings.HasSuffix(fld, "?")
	strippedFieldName := strings.ReplaceAll(fld, "?", "")
	if !regexp.MustCompile(jsVariableNameRegex).Match([]byte(strippedFieldName)) {
		fld = fmt.Sprintf(`"%s"`, fld)
		if isOptional {
			fld += "?"
		}
	}
	if isAnyType {
		fldType = strings.Split(fldType, ".")[0]
		t.fields = append(t.fields, fmt.Sprint(t.indent, "// Go type: ", fldType, "\n", t.indent, fld, ": any;"))
	} else {
		t.fields = append(t.fields, fmt.Sprint(t.indent, fld, ": ", fldType, ";"))
	}
}

func indentLines(str string, i int) string {
	lines := strings.Split(str, "\n")
	for n := range lines {
		lines[n] = strings.Repeat("\t", i) + lines[n]
	}
	return strings.Join(lines, "\n")
}

func getStructFQN(in string) string {
	result := strings.ReplaceAll(in, "[]", "")
	result = strings.ReplaceAll(result, "*", "")
	return result
}

func differentNamespaces(namespace string, typeOf reflect.Type) bool {
	if strings.ContainsRune(typeOf.String(), '.') {
		typeNamespace := strings.Split(typeOf.String(), ".")[0]
		if namespace != typeNamespace {
			return true
		}
	}
	return false
}

func typeClashWithReservedKeyword(input string) bool {
	in := strings.ToLower(strings.TrimSpace(input))
	for _, v := range jsReservedKeywords {
		if in == v {
			return true
		}
	}

	return false
}

func warnAboutTypesClash(entity string) {
	// TODO: Refactor logging
	l := log.New(os.Stderr, "", 0)
	l.Printf("Usage of reserved keyword found and not supported: %s", entity)
	log.Println("Please rename returned type or consider adding bindings config to your wails.json")
}
