package runtime

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

type StoreProvider struct {
	runtime *Runtime
}

func NewStoreProvider(runtime *Runtime) *StoreProvider {
	return &StoreProvider{
		runtime: runtime,
	}
}

// Store is where we keep named data
type Store struct {
	name                string
	data                interface{}
	eventPrefix         string
	callbacks           []func(interface{})
	runtime             *Runtime
	notifySynchronously bool
}

// New creates a new store
func (p *StoreProvider) New(name string) *Store {

	result := Store{
		name:    name,
		runtime: p.runtime,
	}

	result.setupListner()

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

func (s *Store) setupListner() {
	// Setup listener
	s.runtime.Events.On("wails:sync:store:updated:"+s.name, func(data ...interface{}) {

		// store the data
		s.data = data[0]

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

	// Emit event
	s.runtime.Events.Emit("wails:sync:store:updated:"+s.name, s.data)

	// Notify listeners
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
