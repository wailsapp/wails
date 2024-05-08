package collect

import (
	"cmp"
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
	// if a call to [StructInfo.Collect] has completed before the access,
	// for example by calling it in the accessing goroutine
	// or before spawning the accessing goroutine.
	StructInfo struct {
		Fields []*FieldInfo

		typ       *types.Struct
		collector *Collector
		once      sync.Once
	}

	// FieldInfo represents a single field found in a struct.
	FieldInfo struct {
		Field    *types.Var
		Name     string
		Type     types.Type
		Optional bool
		Quoted   bool
	}
)

// Package retrieves the the unique [StructInfo] instance
// associated to the given type within a Collector.
// If none is present, a new one is initialised.
//
// Struct is safe for concurrent use.
func (collector *Collector) Struct(typ *types.Struct) *StructInfo {
	collector.mu.Lock()
	if info := collector.structs.At(typ); info != nil {
		collector.mu.Unlock()
		return info.(*StructInfo)
	}

	info := &StructInfo{
		typ:       typ,
		collector: collector,
	}

	collector.structs.Set(typ, info)
	collector.mu.Unlock()

	return info
}

// Collect gathers information for the structure described by its receiver.
// It can be called concurrently by multiple goroutines;
// the computation will be performed just once.
//
// The field list of the receiver is populated
// by the same flattening algorithm employed by encoding/json.
// JSON struct tags are accounted for.
//
// After Collect returns, the calling goroutine and all goroutines
// it might spawn afterwards are free to access
// the receiver's fields indefinitely.
func (info *StructInfo) Collect() {
	type extField struct {
		*FieldInfo

		// Data for the encoding/json flattening algorithm.
		nameFromTag bool
		index       []int
		info        *StructInfo
	}

	info.once.Do(func() {
		// Flattened list of fields with additional information.
		fields := make([]extField, 0, info.typ.NumFields())

		// Queued embedded types for current and next level.
		current := make([]extField, 0, info.typ.NumFields())
		next := make([]extField, 1, info.typ.NumFields())

		// Count of queued embedded types for current and next level.
		count := make(map[*StructInfo]int)
		nextCount := make(map[*StructInfo]int)

		// Set of visited types to avoid duplicating work.
		visited := make(map[*StructInfo]bool)

		next[0].Type = info.typ

		for len(next) > 0 {
			current, next = next, current[:0]
			count, nextCount = nextCount, count
			clear(nextCount)

			for _, embedded := range current {
				if visited[embedded.info] {
					continue
				}
				visited[embedded.info] = true

				// WARNING: DO NOT EVER CALL einfo.Collect HERE.
				// First, it may deadlock on cyclic types.
				// Second, reusing other structs _after_ flattening
				// may lead to incorrect results for subtle reasons.

				// Scan embedded type for fields to include.
				estruct := embedded.Type.(*types.Struct)
				for i := range estruct.NumFields() {
					field := estruct.Field(i)

					// If the underlying type is a struct, extract it.
					var fstruct *types.Struct
					{
						ftype := types.Unalias(field.Type())
						if ptr, ok := ftype.(*types.Pointer); ok {
							ftype = types.Unalias(ptr.Elem())
						}

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

						finfo := extField{
							FieldInfo: &FieldInfo{
								Field:    field,
								Name:     name,
								Type:     field.Type(),
								Optional: optional,
								Quoted:   quoted,
							},
							nameFromTag: name != "",
							index:       index,
						}

						if name == "" {
							finfo.Name = field.Name()
						}

						fields = append(fields, finfo)
						if count[embedded.info] > 1 {
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
					fsinfo := info.collector.Struct(fstruct)
					nextCount[fsinfo]++
					if nextCount[fsinfo] == 1 {
						next = append(next, extField{
							FieldInfo: &FieldInfo{
								Type: fstruct,
							},
							index: index,
							info:  fsinfo,
						})
					}
				}
			}
		}

		// Prepare for field selection phase.
		slices.SortFunc(fields, func(f1 extField, f2 extField) int {
			// Sort by name first.
			if diff := strings.Compare(f1.Name, f2.Name); diff != 0 {
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
			if j < len(fields) && fields[i].Name == fields[j].Name {
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
		slices.SortFunc(fields, func(f1 extField, f2 extField) int {
			return slices.Compare(f1.index, f2.index)
		})

		// Copy selected fields to receiver.
		info.Fields = make([]*FieldInfo, len(fields))
		for i, field := range fields {
			info.Fields[i] = field.FieldInfo
		}

		info.typ = nil
		info.collector = nil
	})
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
		case "omitempty":
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
