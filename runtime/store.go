// Package runtime contains all the methods and data structures related to the
// runtime library of Wails. This includes both Go and JS runtimes.
package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

// Options defines the optional data that may be used
// when creating a Store
type Options struct {

	// The name of the store
	Name string

	// The runtime to attach the store to
	Runtime *Runtime

	// Indicates if notifying Go listeners should be notified of updates
	// synchronously (on the current thread) or asynchronously using
	// goroutines
	NotifySynchronously bool
}

// StoreProvider is a struct that creates Stores
type StoreProvider struct {
	runtime *Runtime
}

// NewStoreProvider creates new stores using the provided Runtime reference.
func NewStoreProvider(runtime *Runtime) *StoreProvider {
	return &StoreProvider{
		runtime: runtime,
	}
}

// Store is where we keep named data
type Store struct {
	name                string
	data                reflect.Value
	dataType            reflect.Type
	structType          bool
	eventPrefix         string
	callbacks           []func(interface{})
	runtime             *Runtime
	notifySynchronously bool

	// Error handler
	errorHandler func(error)
}

// New creates a new store
func (p *StoreProvider) New(name string, defaultValue interface{}) *Store {

	dataType := reflect.TypeOf(defaultValue)

	result := Store{
		name:       name,
		runtime:    p.runtime,
		data:       reflect.ValueOf(defaultValue),
		dataType:   dataType,
		structType: dataType.Kind() == reflect.Ptr,
	}

	// Setup the sync listener
	result.setupListener()

	return &result
}

// OnError takes a function that will be called
// whenever an error occurs
func (s *Store) OnError(callback func(error)) {
	s.errorHandler = callback
}

// processUpdatedScalar will process the given scalar json
func (s *Store) processUpdatedScalar(data json.RawMessage) error {

	// Unmarshall the value
	var decodedVal interface{}
	err := json.Unmarshal(data, &decodedVal)
	if err != nil {
		return err
	}

	// Convert to correct type
	if decodedVal == nil {
		s.data = reflect.Zero(s.dataType)
	} else {
		s.data = reflect.ValueOf(decodedVal).Convert(s.dataType)
	}

	return nil
}

// processUpdatedStruct will process the given struct json
func (s *Store) processUpdatedStruct(data json.RawMessage) error {

	newData := reflect.New(s.dataType.Elem()).Interface()
	err := json.Unmarshal(data, &newData)
	if err != nil {
		return err
	}
	s.data = reflect.ValueOf(newData)
	return nil
}

// Processes the updates sent by the front end
func (s *Store) processUpdatedData(data string) error {

	var rawdata json.RawMessage
	d := json.NewDecoder(bytes.NewBufferString(data))
	err := d.Decode(&rawdata)
	if err != nil {
		return err
	}

	// If it's a struct process it differently
	if s.structType {
		return s.processUpdatedStruct(rawdata)
	}

	return s.processUpdatedScalar(rawdata)

}

// Setup listener for front end changes
func (s *Store) setupListener() {

	// Listen for updates from the front end
	s.runtime.Events.On("wails:sync:store:updatedbyfrontend:"+s.name, func(data ...interface{}) {

		// Process the incoming data
		err := s.processUpdatedData(data[0].(string))

		if err != nil {
			if s.errorHandler != nil {
				s.errorHandler(err)
				return
			}
		}

		// Notify listeners
		s.notify()
	})
}

// notify the listeners of the current data state
func (s *Store) notify() {

	// Notify callbacks
	for _, callback := range s.callbacks {

		if s.notifySynchronously {
			callback(s.data)
		} else {
			go callback(s.data)
		}

	}
}

// Set will update the data held by the store
// and notify listeners of the change
func (s *Store) Set(data interface{}) error {

	inType := reflect.TypeOf(data)

	if inType != s.dataType {
		return fmt.Errorf("invalid data given in Store.Set(). Expected %s, got %s", s.dataType.String(), inType.String())
	}

	// Save data
	s.data = reflect.ValueOf(data)

	// Stringify data
	newdata, err := json.Marshal(data)
	if err != nil {
		if s.errorHandler != nil {
			return err
		}
	}

	// Emit event to front end
	s.runtime.Events.Emit("wails:sync:store:updatedbybackend:"+s.name, string(newdata))

	// Notify subscribers
	s.notify()

	return nil
}

// Subscribe will subscribe to updates to the store by
// providing a callback. Any updates to the store are sent
// to the callback
func (s *Store) Subscribe(callback func(interface{})) {
	s.callbacks = append(s.callbacks, callback)
}

func (s *Store) updaterCheck(updater interface{}) error {

	// Get type
	updaterType := reflect.TypeOf(updater)

	// Check updater is a function
	if updaterType.Kind() != reflect.Func {
		return fmt.Errorf("invalid value given to store.Update(). Expected 'func(%s) %s'", s.dataType.String(), s.dataType.String())
	}

	// Check input param
	if updaterType.NumIn() != 1 {
		return fmt.Errorf("invalid number of parameters given in updater function. Expected 1")
	}

	// Check input data type
	if updaterType.In(0) != s.dataType {
		return fmt.Errorf("invalid type for input parameter given in updater function. Expected %s, got %s", s.dataType.String(), updaterType.In(0))
	}

	// Check output param
	if updaterType.NumOut() != 1 {
		return fmt.Errorf("invalid number of return parameters given in updater function. Expected 1")
	}

	// Check output data type
	if updaterType.Out(0) != s.dataType {
		return fmt.Errorf("invalid type for return parameter given in updater function. Expected %s, got %s", s.dataType.String(), updaterType.Out(0))
	}

	return nil
}

// Update takes a function that is passed the current state.
// The result of that function is then set as the new state
// of the store. This will notify listeners of the change
func (s *Store) Update(updater interface{}) {

	err := s.updaterCheck(updater)
	if err != nil {
		log.Fatal(err)
	}

	// Build args
	args := []reflect.Value{s.data}

	// Make call
	results := reflect.ValueOf(updater).Call(args)

	// We will only have 1 result. Set the store to it
	s.Set(results[0].Interface())
}
