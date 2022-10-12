package binding

import (
	"encoding/json"
	"sort"
	"sync"
	"unsafe"
)

// DB is our database of method bindings
type DB struct {
	//  map[packagename] -> map[structname] -> map[methodname]*method
	store map[string]map[string]map[string]*BoundMethod

	// This uses fully qualified method names as a shortcut for store traversal.
	// It used for performance gains at runtime
	methodMap map[string]*BoundMethod

	// This uses ids to reference bound methods at runtime
	obfuscatedMethodMap map[int]*BoundMethod

	// Lock to ensure sync access to the data
	lock sync.RWMutex
}

func newDB() *DB {
	return &DB{
		store:               make(map[string]map[string]map[string]*BoundMethod),
		methodMap:           make(map[string]*BoundMethod),
		obfuscatedMethodMap: make(map[int]*BoundMethod),
	}
}

// GetMethodFromStore returns the method for the given package/struct/method names
// nil is returned if any one of those does not exist
func (d *DB) GetMethodFromStore(packageName string, structName string, methodName string) *BoundMethod {

	// Lock the db whilst processing and unlock on return
	d.lock.RLock()
	defer d.lock.RUnlock()

	structMap, exists := d.store[packageName]
	if !exists {
		return nil
	}
	methodMap, exists := structMap[structName]
	if !exists {
		return nil
	}
	return methodMap[methodName]
}

// GetMethod returns the method for the given qualified method name
// qualifiedMethodName is "packagename.structname.methodname"
func (d *DB) GetMethod(qualifiedMethodName string) *BoundMethod {

	// Lock the db whilst processing and unlock on return
	d.lock.RLock()
	defer d.lock.RUnlock()

	return d.methodMap[qualifiedMethodName]
}

// GetObfuscatedMethod returns the method for the given ID
func (d *DB) GetObfuscatedMethod(id int) *BoundMethod {
	// Lock the db whilst processing and unlock on return
	d.lock.RLock()
	defer d.lock.RUnlock()

	return d.obfuscatedMethodMap[id]
}

// AddMethod adds the given method definition to the db using the given qualified path: packageName.structName.methodName
func (d *DB) AddMethod(packageName string, structName string, methodName string, methodDefinition *BoundMethod) {

	// Lock the db whilst processing and unlock on return
	d.lock.Lock()
	defer d.lock.Unlock()

	// Get the map associated with the package name
	structMap, exists := d.store[packageName]
	if !exists {
		// Create a new map for this packagename
		d.store[packageName] = make(map[string]map[string]*BoundMethod)
		structMap = d.store[packageName]
	}

	// Get the map associated with the struct name
	methodMap, exists := structMap[structName]
	if !exists {
		// Create a new map for this packagename
		structMap[structName] = make(map[string]*BoundMethod)
		methodMap = structMap[structName]
	}

	// Store the method definition
	methodMap[methodName] = methodDefinition

	// Store in the methodMap
	key := packageName + "." + structName + "." + methodName
	d.methodMap[key] = methodDefinition

}

// ToJSON converts the method map to JSON
func (d *DB) ToJSON() (string, error) {

	// Lock the db whilst processing and unlock on return
	d.lock.RLock()
	defer d.lock.RUnlock()

	d.UpdateObfuscatedCallMap()

	bytes, err := json.Marshal(&d.store)

	// Return zero copy string as this string will be read only
	result := *(*string)(unsafe.Pointer(&bytes))
	return result, err
}

// UpdateObfuscatedCallMap sets up the secure call mappings
func (d *DB) UpdateObfuscatedCallMap() map[string]int {

	var mappings = make(map[string]int)

	// Iterate map keys and sort them
	keys := make([]string, 0, len(d.methodMap))
	for k := range d.methodMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Iterate sorted keys and add to obfuscated method map
	for id, k := range keys {
		mappings[k] = id
		d.obfuscatedMethodMap[id] = d.methodMap[k]
	}
	return mappings
}
