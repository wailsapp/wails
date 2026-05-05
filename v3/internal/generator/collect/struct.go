package collect

import (
	"cmp"
	"go/ast"
	"go/types"
	"reflect"
	"slices"
	"strings"
	"sync"
	"unicode"
)

type (
	// StructInfo records the flattened field list for a struct type,
	// taking into account JSON tags.
	//
	// The field list is initially empty. It will be populated
	// upon calling [StructInfo.Collect] for the first time.
	//
	// Read accesses to the field list are only safe
	// if a call to [StructInfo.Collect] has been completed before the access,
	// for example by calling it in the accessing goroutine
	// or before spawning the accessing goroutine.
	StructInfo struct {
		Fields []*StructField

		typ *types.Struct

		collector *Collector
		once      sync.Once
	}

	// FieldInfo represents a single field in a struct.
	StructField struct {
		JsonName string // Avoid collisions with [FieldInfo.Name].
		Type     types.Type
		Optional bool
		Quoted   bool

		// Object holds the described type-checker object.
		Object *types.Var
	}
)

func newStructInfo(collector *Collector, typ *types.Struct) *StructInfo {
	return &StructInfo{
		typ:       typ,
		collector: collector,
	}
}

// Struct retrieves the unique [StructInfo] instance
// associated to the given type within a Collector.
// If none is present, a new one is initialised.
//
// Struct is safe for concurrent use.
func (collector *Collector) Struct(typ *types.Struct) *StructInfo {
	// Cache by type pointer, do not use a typeutil.Map:
	//   - for models, it may result in incorrect comments;
	//   - for anonymous structs, it would probably bring little benefit
	//     because the probability of repetitions is much lower.

	return collector.fromCache(typ).(*StructInfo)
}

func (*StructInfo) Object() types.Object {
	return nil
}

func (info *StructInfo) Type() types.Type {
	return info.typ
}

func (*StructInfo) Node() ast.Node {
	return nil
}

// Collect gathers information for the structure described by its receiver.
// It can be called concurrently by multiple goroutines;
// the computation will be performed just once.
//
// The field list of the receiver is populated
// by the same flattening algorithm employed by encoding/json.
// JSON struct tags are accounted for.
//
// Collect returns the receiver for chaining.
// It is safe to call Collect with nil receiver.
//
// After Collect returns, the calling goroutine and all goroutines
// it might spawn afterwards are free to access
// the receiver's fields indefinitely.
func (info *StructInfo) Collect() *StructInfo {
	if info == nil {
		return nil
	}

	type fieldData struct {
		*StructField

		// Data for the encoding/json flattening algorithm.
		nameFromTag bool
		index       []int
	}

	info.once.Do(func() {
		// Flattened list of fields with additional information.
		fields := make([]fieldData, 0, info.typ.NumFields())

		// Queued embedded types for current and next level.
		current := make([]fieldData, 0, info.typ.NumFields())
		next := make([]fieldData, 1, max(1, info.typ.NumFields()))

		// Count of queued embedded types for current and next level.
		count := make(map[*types.Struct]int)
		nextCount := make(map[*types.Struct]int)

		// Set of visited types to avoid duplicating work.
		visited := make(map[*types.Struct]bool)

		next[0] = fieldData{
			StructField: &StructField{
				Type: info.typ,
			},
		}

		for len(next) > 0 {
			current, next = next, current[:0]
			count, nextCount = nextCount, count
			clear(nextCount)

			for _, embedded := range current {
				// Scan embedded type for fields to include.
				estruct := embedded.Type.Underlying().(*types.Struct)

				// Skip previously visited structs
				if visited[estruct] {
					continue
				}
				visited[estruct] = true

				// WARNING: do not reuse cached info for embedded structs.
				// It may lead to incorrect results for subtle reasons.

				for i := range estruct.NumFields() {
					field := estruct.Field(i)

					// Retrieve type of field, following aliases conservatively
					// and unwrapping exactly one pointer.
					ftype := field.Type()
					if ptr, ok := types.Unalias(ftype).(*types.Pointer); ok {
						ftype = ptr.Elem()
					}

					// Detect struct alias and keep it.
					fstruct, _ := types.Unalias(ftype).(*types.Struct)
					if fstruct == nil {
						// Not a struct alias, follow alias chain.
						ftype = types.Unalias(ftype)
						fstruct, _ = ftype.Underlying().(*types.Struct)
					}

					if field.Embedded() {
						if !field.Exported() && fstruct == nil {
							// Ignore embedded fields of unexported non-struct types.
							continue
						}
					} else if !field.Exported() {
						// Ignore unexported non-embedded fields.
						continue
					}

					// Retrieve and parse json tag.
					tag := reflect.StructTag(estruct.Tag(i)).Get("json")
					name, optional, quoted, visible := parseTag(tag)
					if !visible {
						// Ignored by encoding/json.
						continue
					}

					if !isValidFieldName(name) {
						// Ignore alternative name if invalid.
						name = ""
					}

					index := make([]int, len(embedded.index)+1)
					copy(index, embedded.index)
					index[len(embedded.index)] = i

					if name != "" || !field.Embedded() || fstruct == nil {
						// Tag name is non-empty,
						// or field is not embedded,
						// or field is not structure:
						// add to field list.

						if !info.collector.options.UseInterfaces {
							// In class mode, mark parametric fields as optional
							// because there is no way to know their default JS value in advance.
							if _, isTypeParam := types.Unalias(field.Type()).(*types.TypeParam); isTypeParam {
								optional = true
							}
						}

						finfo := fieldData{
							StructField: &StructField{
								JsonName: name,
								Type:     field.Type(),
								Optional: optional,
								Quoted:   quoted,

								Object: field,
							},
							nameFromTag: name != "",
							index:       index,
						}

						if name == "" {
							finfo.JsonName = field.Name()
						}

						fields = append(fields, finfo)
						if count[estruct] > 1 {
							// The struct we are scanning
							// appears multiple times at the current level.
							// This means that all its fields are ambiguous
							// and must disappear.
							// Duplicate them so that the field selection phase
							// below will erase them.
							fields = append(fields, finfo)
						}

						continue
					}

					// Queue embedded field for next level.
					// If it has been queued already, do not duplicate it.
					nextCount[fstruct]++
					if nextCount[fstruct] == 1 {
						next = append(next, fieldData{
							StructField: &StructField{
								Type: ftype,
							},
							index: index,
						})
					}
				}
			}
		}

		// Prepare for field selection phase.
		slices.SortFunc(fields, func(f1 fieldData, f2 fieldData) int {
			// Sort by name first.
			if diff := strings.Compare(f1.JsonName, f2.JsonName); diff != 0 {
				return diff
			}

			// Break ties by depth of occurrence.
			if diff := cmp.Compare(len(f1.index), len(f2.index)); diff != 0 {
				return diff
			}

			// Break ties by presence of json tag (prioritize presence).
			if f1.nameFromTag != f2.nameFromTag {
				if f1.nameFromTag {
					return -1
				} else {
					return 1
				}
			}

			// Break ties by order of occurrence.
			return slices.Compare(f1.index, f2.index)
		})

		fieldCount := 0

		// Keep for each name the dominant field, drop those for which ties
		// still exist (ignoring order of occurrence).
		for i, j := 0, 1; j <= len(fields); j++ {
			if j < len(fields) && fields[i].JsonName == fields[j].JsonName {
				continue
			}

			// If there is only one field with the current name,
			// or there is a dominant one, keep it.
			if i+1 == j || len(fields[i].index) != len(fields[i+1].index) || fields[i].nameFromTag != fields[i+1].nameFromTag {
				fields[fieldCount] = fields[i]
				fieldCount++
			}

			i = j
		}

		fields = fields[:fieldCount]

		// Sort by order of occurrence.
		slices.SortFunc(fields, func(f1 fieldData, f2 fieldData) int {
			return slices.Compare(f1.index, f2.index)
		})

		// Copy selected fields to receiver.
		info.Fields = make([]*StructField, len(fields))
		for i, field := range fields {
			info.Fields[i] = field.StructField
		}

		info.typ = nil
	})

	return info
}

// parseTag parses a json field tag and extracts
// all options recognised by encoding/json.
func parseTag(tag string) (name string, optional bool, quoted bool, visible bool) {
	if tag == "-" {
		return "", false, false, false
	} else {
		visible = true
	}

	parts := strings.Split(tag, ",")

	name = parts[0]

	for _, option := range parts[1:] {
		switch option {
		case "omitempty", "omitzero":
			optional = true
		case "string":
			quoted = true
		}
	}

	return
}

// isValidFieldName determines whether a field name is valid
// according to encoding/json.
func isValidFieldName(name string) bool {
	if name == "" {
		return false
	}

	for _, c := range name {
		if !strings.ContainsRune("!#$%&()*+-./:;<=>?@[]^_{|}~ ", c) && !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return false
		}
	}

	return true
}
