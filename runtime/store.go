// Package runtime contains all the methods and data structures related to the
// runtime library of Wails. This includes both Go and JS runtimes.
package runtime

import (
	"bytes"
	"encoding/json"
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
	data                interface{}
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

	result := Store{
		name:    name,
		runtime: p.runtime,
		data:    defaultValue,
	}

	// initialise the store
	result.init()

	return &result
}

// NewWithOptions creates a new store with the given options
func (p *StoreProvider) NewWithOptions(options Options) *Store {

	result := Store{
		name:                options.Name,
		notifySynchronously: options.NotifySynchronously,
	}

	return &result
}

// OnError takes a function that will be called
// whenever an error occurs
func (s *Store) OnError(callback func(error)) {
	s.errorHandler = callback
}

// init the store
func (s *Store) init() {

	// Get the type of the data
	s.dataType = reflect.TypeOf(s.data)

	// Determine if this is a struct type
	s.structType = s.dataType.Kind() == reflect.Ptr

	// Setup the sync listener
	s.setupListener()
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
		s.data = reflect.Zero(s.dataType).Interface()
	} else {
		s.data = reflect.ValueOf(decodedVal).Convert(s.dataType).Interface()
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
	s.data = newData
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
func (s *Store) Set(data interface{}) {

	// Save data
	s.data = data

	// Stringify data
	newdata, err := json.Marshal(s.data)
	if err != nil {
		if s.errorHandler != nil {
			s.errorHandler(err)
		}
	}

	// Emit event to front end
	s.runtime.Events.Emit("wails:sync:store:updatedbybackend:"+s.name, string(newdata))

	// Notify subscribers
	s.notify()
}

// Subscribe will subscribe to updates to the store by
// providing a callback. Any updates to the store are sent
// to the callback
func (s *Store) Subscribe(callback func(interface{})) {
	s.callbacks = append(s.callbacks, callback)
}

// Update takes a function that is passed the current state.
// The result of that function is then set as the new state
// of the store. This will notify listeners of the change
func (s *Store) Update(updater func(interface{}) interface{}) {
	newData := updater(s.data)
	s.Set(newData)
}
