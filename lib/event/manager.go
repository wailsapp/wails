package event

import (
	"fmt"
	"sync"

	"github.com/wailsapp/wails/lib/interfaces"
	"github.com/wailsapp/wails/lib/logger"
	"github.com/wailsapp/wails/lib/messages"
)

// Manager handles and processes events
type Manager struct {
	incomingEvents chan *messages.EventData
	quitChannel    chan struct{}
	listeners      map[string][]*eventListener
	running        bool
	log            *logger.CustomLogger
	renderer       interfaces.Renderer // Messages will be dispatched to the frontend
	wg             sync.WaitGroup
	mu             sync.Mutex
}

// NewManager creates a new event manager with a 100 event buffer
func NewManager() interfaces.EventManager {
	return &Manager{
		incomingEvents: make(chan *messages.EventData, 100),
		quitChannel:    make(chan struct{}, 1),
		listeners:      make(map[string][]*eventListener),
		running:        false,
		log:            logger.NewCustomLogger("Events"),
	}
}

// PushEvent places the given event on to the event queue
func (e *Manager) PushEvent(eventData *messages.EventData) {
	e.incomingEvents <- eventData
}

// eventListener holds a callback function which is invoked when
// the event listened for is emitted. It has a counter which indicates
// how the total number of events it is interested in. A value of zero
// means it does not expire (default).
type eventListener struct {
	callback func(...interface{}) // Function to call with emitted event data
	counter  uint                 // Expire after counter callbacks. 0 = infinite
	expired  bool                 // Indicates if the listener has expired
}

// Creates a new event listener from the given callback function
func (e *Manager) addEventListener(eventName string, callback func(...interface{}), counter uint) error {

	// Sanity check inputs
	if callback == nil {
		return fmt.Errorf("nil callback bassed to addEventListener")
	}

	// Check event has been registered before
	if e.listeners[eventName] == nil {
		e.listeners[eventName] = []*eventListener{}
	}

	// Create the callback
	listener := &eventListener{
		callback: callback,
		counter:  counter,
	}

	// Register listener
	e.listeners[eventName] = append(e.listeners[eventName], listener)

	// All good mate
	return nil
}

// On adds a listener for the given event
func (e *Manager) On(eventName string, callback func(...interface{})) {
	// Add a persistent eventListener (counter = 0)
	err := e.addEventListener(eventName, callback, 0)
	if err != nil {
		e.log.Error(err.Error())
	}
}

// Once adds a listener for the given event that will auto remove
// after one callback
func (e *Manager) Once(eventName string, callback func(...interface{})) {
	// Add a persistent eventListener (counter = 0)
	err := e.addEventListener(eventName, callback, 1)
	if err != nil {
		e.log.Error(err.Error())
	}
}

// OnMultiple adds a listener for the given event that will trigger
// at most <counter> times.
func (e *Manager) OnMultiple(eventName string, callback func(...interface{}), counter uint) {
	// Add a persistent eventListener (counter = 0)
	err := e.addEventListener(eventName, callback, counter)
	if err != nil {
		e.log.Error(err.Error())
	}
}

// Emit broadcasts the given event to the subscribed listeners
func (e *Manager) Emit(eventName string, optionalData ...interface{}) {
	e.incomingEvents <- &messages.EventData{Name: eventName, Data: optionalData}
}

// Start the event manager's queue processing
func (e *Manager) Start(renderer interfaces.Renderer) {

	e.log.Info("Starting")

	// Store renderer
	e.renderer = renderer

	// Set up waitgroup so we can wait for goroutine to quit
	e.running = true
	e.wg.Add(1)

	// Run main loop in separate goroutine
	go func() {
		e.log.Info("Listening")
		for e.running {
			// TODO: Listen for application exit
			select {
			case event := <-e.incomingEvents:
				e.log.DebugFields("Got Event", logger.Fields{
					"data": event.Data,
					"name": event.Name,
				})

				// Notify renderer
				err := e.renderer.NotifyEvent(event)
				if err != nil {
					e.log.Error(err.Error())
				}

				e.mu.Lock()

				// Iterate listeners
				for _, listener := range e.listeners[event.Name] {

					if !listener.expired {
						// Call listener, perhaps with data
						if event.Data == nil {
							go listener.callback()
						} else {
							unpacked := event.Data.([]interface{})
							go listener.callback(unpacked...)
						}
					}

					// Update listen counter
					if listener.counter > 0 {
						listener.counter = listener.counter - 1
						if listener.counter == 0 {
							listener.expired = true
						}
					}
				}

				e.mu.Unlock()

			case <-e.quitChannel:
				e.running = false
			}
		}
		e.wg.Done()
	}()
}

// Shutdown is called when exiting the Application
func (e *Manager) Shutdown() {
	e.log.Debug("Shutting Down")
	e.quitChannel <- struct{}{}
	e.log.Debug("Waiting for main loop to exit")
	e.wg.Wait()
}
