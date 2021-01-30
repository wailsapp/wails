// package runtime contains all the methods and data structures related to the
// runtime library of Wails. This includes both Go and JS runtimes.
package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	golog "log"
	"os"
	"reflect"
	"sync"

	"github.com/wailsapp/wails/v2/internal/deepcopy"
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

// newStore creates new stores using the provided Runtime reference.
func newStore(runtime *Runtime) *StoreProvider {
	return &StoreProvider{
		runtime: runtime,
	}
}

// Store is where we keep named data
type Store struct {
	name                string
	data                reflect.Value
	dataType            reflect.Type
	eventPrefix         string
	callbacks           []reflect.Value
	runtime             *Runtime
	notifySynchronously bool

	// Lock
	mux sync.Mutex

	// Error handler
	errorHandler func(error)
}

func fatal(err error) {
	println(err.Error())
	os.Exit(1)
}

// New creates a new store
func (p *StoreProvider) New(name string, defaultValue interface{}) *Store {

	if defaultValue == nil {
		golog.Fatal("Cannot initialise a store with nil")
	}

	result := Store{
		name:    name,
		runtime: p.runtime,
	}

	// Setup the sync listener
	result.setupListener()

	result.Set(defaultValue)

	return &result
}

func (s *Store) lock() {
	s.mux.Lock()
}

func (s *Store) unlock() {
	s.mux.Unlock()
}

// OnError takes a function that will be called
// whenever an error occurs
func (s *Store) OnError(callback func(error)) {
	s.errorHandler = callback
}

// Processes the updates sent by the front end
func (s *Store) processUpdatedData(data string) error {

	// Decode incoming data
	var rawdata json.RawMessage
	d := json.NewDecoder(bytes.NewBufferString(data))
	err := d.Decode(&rawdata)
	if err != nil {
		return err
	}

	// Create a new instance of our data and unmarshal
	// the received value into it
	newData := reflect.New(s.dataType).Interface()
	err = json.Unmarshal(rawdata, &newData)
	if err != nil {
		return err
	}

	// Lock mutex for writing
	s.lock()

	// Handle nulls
	if newData == nil {
		s.data = reflect.Zero(s.dataType)
	} else {
		// Store the resultant value in the data store
		s.data = reflect.ValueOf(newData).Elem()
	}

	// Unlock mutex
	s.unlock()

	return nil
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

	// Listen for resync events
	s.runtime.Events.On("wails:sync:store:resync:"+s.name, func(data ...interface{}) {
		// Resetting the curent data will resync
		s.resync()
	})

	// Do initial resync
	s.resync()
}

func (s *Store) resync() {

	// Lock
	s.lock()
	defer s.unlock()

	var result string

	if s.data.IsValid() {
		rawdata, err := json.Marshal(s.data.Interface())
		if err != nil {
			if s.errorHandler != nil {
				s.errorHandler(err)
				return
			}
		}
		result = string(rawdata)
	} else {
		result = "{}"
	}

	// Emit event to front end
	s.runtime.Events.Emit("wails:sync:store:updatedbybackend:"+s.name, result)

	// Notify subscribers
	s.notify()
}

// notify the listeners of the current data state
func (s *Store) notify() {

	// Execute callbacks
	for _, callback := range s.callbacks {

		// Build args
		s.lock()
		args := []reflect.Value{s.data}
		s.unlock()

		if s.notifySynchronously {
			callback.Call(args)
		} else {
			go callback.Call(args)
		}

	}
}

// Set will update the data held by the store
// and notify listeners of the change
func (s *Store) Set(data interface{}) error {

	if data == nil {
		return fmt.Errorf("cannot set store to nil")
	}

	s.lock()

	dataCopy := deepcopy.Copy(data)

	if dataCopy != nil {
		inType := reflect.TypeOf(dataCopy)

		if inType != s.dataType && s.data.IsValid() {
			s.unlock()
			return fmt.Errorf("invalid data given in Store.Set(). Expected %s, got %s", s.dataType.String(), inType.String())
		}
	}

	if s.dataType == nil {
		s.dataType = reflect.TypeOf(dataCopy)
	}

	// Save data
	s.data = reflect.ValueOf(dataCopy)

	s.unlock()

	// Resync with subscribers
	s.resync()

	return nil
}

// callbackCheck ensures the given function to Subscribe() is
// of the correct signature. Absolutely cannot wait for
// generics to land rather than writing this nonsense.
func (s *Store) callbackCheck(callback interface{}) error {

	// Get type
	callbackType := reflect.TypeOf(callback)

	// Check callback is a function
	if callbackType.Kind() != reflect.Func {
		return fmt.Errorf("invalid value given to store.Subscribe(). Expected 'func(%s)'", s.dataType.String())
	}

	// Check input param
	if callbackType.NumIn() != 1 {
		return fmt.Errorf("invalid number of parameters given in callback function. Expected 1")
	}

	// Check input data type
	if callbackType.In(0) != s.dataType {
		return fmt.Errorf("invalid type for input parameter given in callback function. Expected %s, got %s", s.dataType.String(), callbackType.In(0))
	}

	// Check output param
	if callbackType.NumOut() != 0 {
		return fmt.Errorf("invalid number of return parameters given in callback function. Expected 0")
	}

	return nil
}

// Subscribe will subscribe to updates to the store by
// providing a callback. Any updates to the store are sent
// to the callback
func (s *Store) Subscribe(callback interface{}) {

	err := s.callbackCheck(callback)
	if err != nil {
		fatal(err)
	}

	callbackFunc := reflect.ValueOf(callback)

	s.lock()
	s.callbacks = append(s.callbacks, callbackFunc)
	s.unlock()
}

// updaterCheck ensures the given function to Update() is
// of the correct signature. Absolutely cannot wait for
// generics to land rather than writing this nonsense.
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
		fatal(err)
	}

	// Build args
	s.lock()
	args := []reflect.Value{s.data}
	s.unlock()

	// Make call
	results := reflect.ValueOf(updater).Call(args)

	// We will only have 1 result. Set the store to it
	s.Set(results[0].Interface())
}

// Get returns the value of the data that's kept in the current state / Store
func (s *Store) Get() interface{} {
	s.lock()
	defer s.unlock()

	if !s.data.IsValid() {
		return nil
	}

	return s.data.Interface()
}
